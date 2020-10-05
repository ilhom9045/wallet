package wallet

import (
	"errors"

	"github.com/google/uuid"
	"github.com/ilhom9045/wallet/pkg/types"
)

//Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
}

//ErrAccountNotFound ...
var ErrAccountNotFound = errors.New("Account ID not found")

//ErrPhoneRegistered ...
var ErrPhoneRegistered = errors.New("This phone was registerred")

//ErrPaymentNotFound ...
var ErrPaymentNotFound = errors.New("Payment ID not found")

//ErrAmountMustBePositive ...
var ErrAmountMustBePositive = errors.New("Amount mast be <0")

//ErrNotEnoughBalance ...
var ErrNotEnoughBalance = errors.New("")

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

//FindPaymentByID ...
func (s *Service) FindPaymentByID(paymentID string) (*types.Payment, error) {
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			return payment, nil
		}
	}
	return nil, ErrPaymentNotFound
}

//Deposit ...
func (s *Service) Deposit(accountID int64, amount types.Money) error {
	if amount < 0 {
		return ErrAmountMustBePositive
	}
	account, err := s.FindAccountByID(accountID)
	if err != nil {
		return err
	}
	account.Balance += amount
	return nil
}

//Pay ...
func (s *Service) Pay(accountID int64, amount types.Money, category types.PaymentCategory) (*types.Payment, error) {
	if amount <= 0 {
		return nil, ErrAmountMustBePositive
	}
	var account *types.Account
	for _, acc := range s.accounts {
		if acc.ID == accountID {
			account = acc
			break
		}
	}
	if account == nil {
		return nil, ErrAccountNotFound
	}
	if account.Balance < amount {
		return nil, ErrNotEnoughBalance
	}
	account.Balance -= amount
	paymentID := uuid.New().String()
	payment := &types.Payment{
		ID:        paymentID,
		Amount:    amount,
		AccountID: accountID,
		Category:  category,
		Status:    types.PaymentStatusInProgress,
	}
	s.payments = append(s.payments, payment)
	return payment, nil

}

//Reject ...
func (s *Service) Reject(paymentID string) error {
	var targetPay *types.Payment
	for _, payment := range s.payments {
		if payment.ID == paymentID {
			targetPay = payment
			break
		}
	}
	if targetPay == nil {
		return ErrAccountNotFound
	}

	acc, err := s.FindAccountByID(targetPay.AccountID)
	if err != nil {
		return err
	}

	targetPay.Status = types.PaymentStatusFail
	acc.Balance += targetPay.Amount
	return nil
}

//Repeat ...
func (s *Service) Repeat(paymentID string) (*types.Payment, error) {
	payment, err := s.FindPaymentByID(paymentID)
	if err != nil {
		return nil, err
	}
	pay, err := s.Pay(payment.AccountID, payment.Amount, payment.Category)
	if err != nil {
		return nil, err
	}
	return pay, nil
}
