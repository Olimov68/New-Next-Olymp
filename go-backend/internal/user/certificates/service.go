package usercertificates

import (
	"errors"

	"gorm.io/gorm"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetMyCertificates(userID uint) ([]CertificateResponse, error) {
	list, err := s.repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]CertificateResponse, len(list))
	for i, c := range list {
		result[i] = ToCertificateResponse(&c)
	}
	return result, nil
}

func (s *Service) GetByID(id, userID uint) (*CertificateResponse, error) {
	c, err := s.repo.GetByID(id, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("certificate not found")
		}
		return nil, err
	}
	res := ToCertificateResponse(c)
	return &res, nil
}
