package logic

import "ledger/internal/adapters"

type service struct {
	repo adapters.Repository
}

func NewService(r adapters.Repository) (Service, error) {
	return &service{
		repo: r,
	}, nil
}
