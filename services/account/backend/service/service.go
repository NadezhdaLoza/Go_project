package service

import "github.com/nadezhda/documents/services/account/domain"

type Service interface {
	CreateAccount(account string, bankName string, bankIdentityCode string) (*domain.Account, error)
	GetAccounts() ([]domain.Account, error)
	UpdateAccount(domain.Account) (*domain.Account, error)
	DeleteAccount(int) error
	AccountById(int) (*domain.Account, error)
}
