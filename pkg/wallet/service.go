package wallet

import (
	"errors"

	"github.com/ilhom9045/wallet/pkg/types"
)

//Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
}

//ErrAccountNotFound ...
var ErrAccountNotFound = errors.New("Account ID not found")

//ErrPhoneRegistered ...
var ErrPhoneRegistered = errors.New("This phone was registerred")

//RegisterAccount ...
func (s *Service) RegisterAccount(phone types.Phone) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.Phone == phone {
			return nil, ErrPhoneRegistered
		}
	}
	s.nextAccountID++
	account := &types.Account{
		ID:      s.nextAccountID,
		Phone:   phone,
		Balance: 0,
	}
	s.accounts = append(s.accounts, account)
	return account, nil
}

//FindAccountByID ...
func (s *Service) FindAccountByID(accountID int64) (*types.Account, error) {
	for _, account := range s.accounts {
		if account.ID == accountID {
			return account, nil
		}
	}
	return nil, ErrAccountNotFound
}
