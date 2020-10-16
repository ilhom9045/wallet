package wallet

import (
	"bufio"
	"errors"
	"github.com/google/uuid"
	"github.com/ilhom9045/wallet/pkg/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

//Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
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
var ErrFavoriteNotFound = errors.New("Favorite not found")
var ErrFileNotFound = errors.New("File not found")

//AccountsFile some text
var AccountsFile = "/accounts.dump"

//some text
var PaymentsFile = "/payments.dump"

//some text
var FavoritesFile = "/favorites.dump"

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
func (s *Service) FavoritePayment(paymentID string, name string) (*types.Favorite, error) {
	var tarPay *types.Payment
	for _, tarp := range s.payments {
		if tarp.ID == paymentID {
			tarPay = tarp
			break
		}
	}
	if tarPay == nil {
		return nil, ErrPaymentNotFound
	}
	pay, err := s.Pay(tarPay.AccountID, tarPay.Amount, tarPay.Category)
	if err != nil {
		return nil, err
	}

	favorite := &types.Favorite{
		ID:        pay.ID,
		AccountID: pay.AccountID,
		Name:      name,
		Amount:    pay.Amount,
		Category:  pay.Category,
	}
	s.favorites = append(s.favorites, favorite)
	return favorite, nil

}
func (s *Service) PayFromFavorite(favoriteID string) (*types.Payment, error) {
	var tarPay *types.Favorite
	for _, favorite := range s.favorites {
		if favorite.ID == favoriteID {
			tarPay = favorite
		}
	}
	if tarPay == nil {
		return nil, ErrFavoriteNotFound
	}
	pay, err := s.Pay(tarPay.AccountID, tarPay.Amount, tarPay.Category)
	if err != nil {
		return nil, err
	}
	return pay, nil
}

var ErrFileNotClose = errors.New("File not close")

func (s *Service) GetDir() (string, error) {
	dir, err := filepath.Abs(".")
	if err != nil {
		return "", err
	}
	return dir, nil
}
func (s *Service) ExportToFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		log.Print("err open file")
		return ErrFileNotFound
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Print(ErrFileNotClose)
		}
	}()
	data := ""
	for _, account := range s.accounts {
		data += strconv.Itoa(int(account.ID)) + ";"
		data += string(account.Phone) + ";"
		data += strconv.Itoa(int(account.Balance)) + "|"
	}

	_, err = file.Write([]byte(data))
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	return nil
}
func (s *Service) ImportFromFile(path string) error {
	file, err := os.Open(path)

	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	//defer closeFile(file)
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	//log.Printf("%#v", file)

	content := make([]byte, 0)
	buf := make([]byte, 4)
	for {
		read, err := file.Read(buf)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
		content = append(content, buf[:read]...)
	}

	data := string(content)

	accounts := strings.Split(data, "|")
	accounts = accounts[:len(accounts)-1]
	// if accounts == nil {
	// 	return ErrAccountNotFound
	// }
	for _, account := range accounts {

		value := strings.Split(account, ";")
		id, err := strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		phone := types.Phone(value[1])
		balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		editAccount := &types.Account{
			ID:      int64(id),
			Phone:   phone,
			Balance: types.Money(balance),
		}

		s.accounts = append(s.accounts, editAccount)
		log.Print(account)
	}
	return nil
}
func (s *Service) Export(dir string) error {
	err := s.exportAccount(dir + AccountsFile)
	if err != nil {
		return err
	}
	err = s.exportPayment(dir + PaymentsFile)
	if err != nil {
		return err
	}
	err = s.exportFavorite(dir + FavoritesFile)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) exportAccount(path string) error {
	if s.accounts == nil {
		log.Println("s.account == nil")
		return nil
	}
	file, err := os.Create(path)
	if err != nil {
		log.Print("err open file")
		return ErrFileNotFound
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Print(err)
		}
	}()
	data := ""
	for _, account := range s.accounts {
		data += strconv.Itoa(int(account.ID)) + ", "
		data += string(account.Phone) + ", "
		data += strconv.Itoa(int(account.Balance)) + "\n"
	}
	_, err = file.Write([]byte(data))
	if err != nil {
		log.Print(err)
		return ErrFileNotFound
	}
	return nil
}
func (s *Service) exportFavorite(path string) error {
	if s.favorites == nil {
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}

	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	favoriteData := ""
	for _, favorite := range s.favorites {
		favoriteData += favorite.ID + ", "
		favoriteData += strconv.Itoa(int(favorite.AccountID)) + ", "
		favoriteData += favorite.Name + ", "
		favoriteData += strconv.Itoa(int(favorite.Amount)) + ", "
		favoriteData += string(favorite.Category) + "\n"
	}
	_, err = file.WriteString(favoriteData)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) exportPayment(path string) error {
	if s.payments == nil {
		return nil
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	paymentData := ""

	for _, payment := range s.payments {
		paymentData += payment.ID + ", "
		paymentData += strconv.Itoa(int(payment.AccountID)) + ", "
		paymentData += strconv.Itoa(int(payment.Amount)) + ", "
		paymentData += string(payment.Category) + ", "
		paymentData += string(payment.Status) + "\n"
	}

	_, err = file.WriteString(paymentData)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) Import(dir string) error {
	err := s.importAccount(dir + AccountsFile)
	if err != nil {
		return err
	}
	err = s.importPayment(dir + PaymentsFile)
	if err != nil {
		return err
	}
	err = s.importFavorite(dir + FavoritesFile)
	if err != nil {
		return err
	}
	return nil
}
func (s *Service) importAccount(path string) error {
	lines, err := readLine(path)
	if err != nil {
		return err
	}
	value := []string{}
	newAccount := &types.Account{}
	for _, line := range lines {
		value = strings.Split(line, ", ")
		ID, err := strconv.Atoi(value[0])
		if err != nil {
			return err
		}
		Balance, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		newAccount.ID = int64(ID)
		newAccount.Phone = types.Phone(value[1])
		newAccount.Balance = types.Money(Balance)
		s.accounts = append(s.accounts, newAccount)
		log.Print(line)
	}
	return nil
}
func (s *Service) importPayment(path string) error {
	lines, err := readLine(path)
	if err != nil {
		return err
	}
	value := []string{}
	newPayment := &types.Payment{}
	for _, line := range lines {
		value = strings.Split(line, ", ")
		ID := value[0]

		AccountID, err := strconv.Atoi(value[1])
		if err != nil {
			return err
		}
		Amount, err := strconv.Atoi(value[2])
		if err != nil {
			return err
		}
		Category := value[3]
		Status := value[4]
		newPayment.ID = ID
		newPayment.AccountID = int64(AccountID)
		newPayment.Amount = types.Money(Amount)
		newPayment.Category = types.PaymentCategory(Category)
		newPayment.Status = types.PaymentStatus(Status)
		s.payments = append(s.payments, newPayment)
		log.Print(line)
	}
	return nil
}
func (s *Service) importFavorite(path string) error {
	lines, err := readLine(path)
	if err != nil {
		return err
	}
	value := []string{}
	newFavorite := &types.Favorite{}
	for _, line := range lines {
		value = strings.Split(line, ", ")
		ID := value[0]

		AccountID, err := strconv.Atoi(value[1])
		if err != nil {
			return err
		}
		Name := value[2]
		Amount, err := strconv.Atoi(value[3])
		if err != nil {
			return err
		}
		Category := value[4]
		newFavorite.ID = ID
		newFavorite.AccountID = int64(AccountID)
		newFavorite.Name = Name
		newFavorite.Amount = types.Money(Amount)
		newFavorite.Category = types.PaymentCategory(Category)
		s.favorites = append(s.favorites, newFavorite)
		log.Print(line)
	}
	return nil
}
func readLine(path string) (lines []string, err error) {

	file, err := os.Open(path)
	if err != nil {
		log.Print(err)
		return nil, ErrFileNotFound
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err != nil || len(line) == 0 {
			break
		}
		lines = append(lines, string(line))
	}
	return lines, nil
}
