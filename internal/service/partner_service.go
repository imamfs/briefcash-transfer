package service

import (
	"briefcash-transfer/internal/entity"
	"briefcash-transfer/internal/helper/loghelper"
	"briefcash-transfer/internal/repository"
	"context"
	"sync"

	"github.com/sirupsen/logrus"
)

const DefaultBankConfig = "14"

type BankPartner interface {
	LoadAllBankPartner(ctx context.Context) error
	GetBankConfig(bankCode string, log *logrus.Entry) entity.BankConfig
}

type bankPartner struct {
	rwMutex     sync.RWMutex
	partnerRepo repository.PartnerRepository
	bankCache   map[string]entity.BankConfig
}

func NewPartnerService(partnerRepo repository.PartnerRepository) BankPartner {
	return &bankPartner{
		partnerRepo: partnerRepo,
	}
}

func (s *bankPartner) LoadAllBankPartner(ctx context.Context) error {
	log := loghelper.Logger.WithFields(logrus.Fields{
		"service":   "partner_service",
		"operation": "load_config",
	})

	// Get list bank partner config from database
	log.Info("Collect bank partner config from database")
	banks, err := s.partnerRepo.FindAll(ctx)
	if err != nil {
		log.WithError(err).Error("Failed to collect bank config from database")
		return err
	}

	// Write bank partner config to memory
	log.Infof("Cache bank partner config to memory, with total data %d", len(banks))
	s.rwMutex.Lock()
	s.bankCache = make(map[string]entity.BankConfig)
	for _, bank := range banks {
		s.bankCache[bank.BankCode] = bank
	}
	s.rwMutex.Unlock()
	return nil
}

func (s *bankPartner) GetBankConfig(bankCode string, log *logrus.Entry) entity.BankConfig {
	// get bank config from memory based on bank bank bankCode
	log.Info("Get bank partner config from memory")
	s.rwMutex.RLock()
	bank, ok := s.bankCache[bankCode]
	s.rwMutex.Unlock()

	// if data config not found, fallback to default bank
	if !ok {
		log.Info("Bank partner not found, fallback to default bank")
		s.rwMutex.RLock()
		fallback := s.bankCache[DefaultBankConfig]
		s.rwMutex.RUnlock()
		log.Infof("Bank %s selected", bank.BankName)
		return fallback
	}

	log.Infof("Bank %s selected", bank.BankName)
	return bank
}
