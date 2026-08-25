package chem

import (
	"labvent/internal/store"
)

type Rule struct {
	From        string  `json:"from"`
	To          string  `json:"to"`
	MinQuantity float64 `json:"min_quantity"`
}

type ClassifyService struct {
	blobs store.Blob
}

func NewClassifyService(blobs store.Blob) *ClassifyService {
	return &ClassifyService{blobs: blobs}
}

func (s *ClassifyService) Register(chemID string, class string) error {
	return s.blobs.Save("chem-class", chemID, ClassRecord{ChemID: chemID, Class: class})
}

func (s *ClassifyService) SetRules(rules []Rule) error {
	return s.blobs.Save("class-rule", "active", RuleSet{Rules: rules})
}

func (s *ClassifyService) Rules() ([]Rule, error) {
	var set RuleSet
	if err := s.blobs.Load("class-rule", "active", &set); err != nil {
		return nil, err
	}
	return set.Rules, nil
}

func (s *ClassifyService) ClassOf(chemID string) (string, error) {
	var record ClassRecord
	if err := s.blobs.Load("chem-class", chemID, &record); err != nil {
		return "", err
	}
	var chem Chemical
	if err := s.blobs.Load("chem", chemID, &chem); err != nil {
		return "", err
	}
	rules, err := s.Rules()
	if err != nil {
		return "", err
	}
	return applyRules(record.Class, chem.Quantity, rules), nil
}

func (s *ClassifyService) Refresh(chemIDs []string) error {
	for _, chemID := range chemIDs {
		class, err := s.ClassOf(chemID)
		if err != nil {
			return err
		}
		if err := s.blobs.Save("chem-class", chemID, ClassRecord{ChemID: chemID, Class: class}); err != nil {
			return err
		}
		var chem Chemical
		if err := s.blobs.Load("chem", chemID, &chem); err != nil {
			return err
		}
		chem.Class = class
		if err := s.blobs.Save("chem", chemID, chem); err != nil {
			return err
		}
	}
	return nil
}

func (s *ClassifyService) RefreshAll() error {
	ids, err := s.blobs.List("chem")
	if err != nil {
		return err
	}
	return s.Refresh(ids)
}

func applyRules(baseClass string, quantity float64, rules []Rule) string {
	for _, rule := range rules {
		if baseClass == rule.From && quantity >= rule.MinQuantity {
			return rule.To
		}
	}
	return baseClass
}

type RuleSet struct {
	Rules []Rule `json:"rules"`
}

type ClassRecord struct {
	ChemID string `json:"chem_id"`
	Class  string `json:"class"`
}
