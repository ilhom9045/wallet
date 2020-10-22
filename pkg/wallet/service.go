package wallet

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/ilhom9045/wallet/pkg/types"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//Service ...
type Service struct {
	nextAccountID int64
	accounts      []*types.Account
	payments      []*types.Payment
	favorites     []*types.Favorite
}
type Progress struct {
	Part   int
	Result types.Money
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
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Print(cerr)
		}
	}()
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
	lenAccounts := len(s.accounts)

	if lenAccounts != 0 {
		fileDir := dir + "/accounts.dump"
		file, err := os.Create(fileDir)
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		data := ""
		for _, account := range s.accounts {
			id := strconv.Itoa(int(account.ID)) + ";"
			phone := string(account.Phone) + ";"
			balance := strconv.Itoa(int(account.Balance))

			data += id
			data += phone
			data += balance + "|"
		}

		_, err = file.Write([]byte(data))
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
	}

	lenPayments := len(s.payments)

	if lenPayments != 0 {
		fileDir := dir + "/payments.dump"

		file, err := os.Create(fileDir)
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		data := ""
		for _, payment := range s.payments {
			idPayment := string(payment.ID) + ";"
			idPaymnetAccountId := strconv.Itoa(int(payment.AccountID)) + ";"
			amountPayment := strconv.Itoa(int(payment.Amount)) + ";"
			categoryPayment := string(payment.Category) + ";"
			statusPayment := string(payment.Status)

			data += idPayment
			data += idPaymnetAccountId
			data += amountPayment
			data += categoryPayment
			data += statusPayment + "|"
		}

		_, err = file.Write([]byte(data))
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
	}

	lenFavorites := len(s.favorites)

	if lenFavorites != 0 {
		fileDir := dir + "/favorites.dump"
		file, err := os.Create(fileDir)
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}

		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()
		data := ""
		for _, favorite := range s.favorites {
			idFavorite := string(favorite.ID) + ";"
			idFavoriteAccountId := strconv.Itoa(int(favorite.AccountID)) + ";"
			nameFavorite := string(favorite.Name) + ";"
			amountFavorite := strconv.Itoa(int(favorite.Amount)) + ";"
			categoryFavorite := string(favorite.Category)

			data += idFavorite
			data += idFavoriteAccountId
			data += nameFavorite
			data += amountFavorite
			data += categoryFavorite + "|"
		}
		_, err = file.Write([]byte(data))
		if err != nil {
			log.Print(err)
			return ErrFileNotFound
		}
	}
	return nil
}
func (s *Service) Import(dir string) error {
	dirAccount := dir + "/accounts.dump"
	file, err := os.Open(dirAccount)
	if err != nil {
		log.Print(err)
		// return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			if cerr := file.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		content := make([]byte, 0)
		buf := make([]byte, 4)
		for {
			read, err := file.Read(buf)
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Print(err)
				//log.Print(dirAccount, " 3333")
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
			//log.Print(editAccount, " read")

			s.accounts = append(s.accounts, editAccount)
		}
	}

	dirPaymnet := dir + "/payments.dump"
	filePayment, err := os.Open(dirPaymnet)

	if err != nil {
		log.Print(err)
		// return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			if cerr := filePayment.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		contentPayment := make([]byte, 0)
		buf := make([]byte, 4)
		for {
			readPayment, err := filePayment.Read(buf)
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Print(err)
				return ErrFileNotFound
			}
			contentPayment = append(contentPayment, buf[:readPayment]...)
		}

		data := string(contentPayment)

		payments := strings.Split(data, "|")
		payments = payments[:len(payments)-1]
		//log.Print(favorites, " fav")
		for _, payment := range payments {

			value := strings.Split(payment, ";")
			idPayment := string(value[0])

			accountIdPeyment, err := strconv.Atoi(value[1])
			if err != nil {
				return err
			}

			amountPayment, err := strconv.Atoi(value[2])
			if err != nil {
				return err
			}
			categoryPayment := types.PaymentCategory(value[3])

			statusPayment := types.PaymentStatus(value[4])
			newPayment := &types.Payment{
				ID:        idPayment,
				AccountID: int64(accountIdPeyment),
				Amount:    types.Money(amountPayment),
				Category:  categoryPayment,
				Status:    statusPayment,
			}

			s.payments = append(s.payments, newPayment)
			//log.Print(payment)

		}
	}

	dirfavorite := dir + "/favorites.dump"
	fileFavorite, err := os.Open(dirfavorite)

	if err != nil {
		log.Print(err)
		// return ErrFileNotFound
		err = ErrFileNotFound
	}
	if err != ErrFileNotFound {
		defer func() {
			if cerr := fileFavorite.Close(); cerr != nil {
				log.Print(cerr)
			}
		}()

		contentFavorite := make([]byte, 0)
		buf := make([]byte, 4)
		for {
			readFavorite, err := fileFavorite.Read(buf)
			if err == io.EOF {
				break
			}

			if err != nil {
				log.Print(err)
				return ErrFileNotFound
			}
			contentFavorite = append(contentFavorite, buf[:readFavorite]...)
		}

		data := string(contentFavorite)
		//log.Print(dirfavorite, " fav ", data)
		favorites := strings.Split(data, "|")
		favorites = favorites[:len(favorites)-1]

		for _, favorite := range favorites {

			valueFavorite := strings.Split(favorite, ";")
			idFavorite := string(valueFavorite[0])
			accountIdFavorite, err := strconv.Atoi(valueFavorite[1])
			if err != nil {
				return err
			}
			nameFavorite := string(valueFavorite[2])

			amountFavorite, err := strconv.Atoi(valueFavorite[3])
			if err != nil {
				return err
			}
			categoryPayment := types.PaymentCategory(valueFavorite[4])

			newFavorite := &types.Favorite{
				ID:        idFavorite,
				AccountID: int64(accountIdFavorite),
				Name:      nameFavorite,
				Amount:    types.Money(amountFavorite),
				Category:  categoryPayment,
			}

			s.favorites = append(s.favorites, newFavorite)
			//log.Print(favorite)
		}
	}
	return nil
}
func (s *Service) ExportAccountHistory(accountID int64) (newPayment []types.Payment, err error) {
	for _, value := range s.payments {
		if value.AccountID == accountID {
			newPayment = append(newPayment, *value)
		}
	}
	if newPayment == nil {
		return nil, ErrAccountNotFound
	}
	return
}
func (s *Service) HistoryToFiles(payments []types.Payment, dir string, records int) error {
	if len(payments) > 0 {
		if len(payments) <= records {
			file, _ := os.OpenFile(dir+"/payments.dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			defer func() {
				if cerr := file.Close(); cerr != nil {
					log.Print(cerr)
				}
			}()

			var str string
			for _, payment := range payments {
				idPayment := payment.ID + ";"
				idPaymnetAccountId := strconv.Itoa(int(payment.AccountID)) + ";"
				amountPayment := strconv.Itoa(int(payment.Amount)) + ";"
				categoryPayment := string(payment.Category) + ";"
				statusPayment := string(payment.Status)

				str += idPayment
				str += idPaymnetAccountId
				str += amountPayment
				str += categoryPayment
				str += statusPayment + "\n"
			}
			_, err := file.WriteString(str)
			if err != nil {
				log.Print(err)
			}
		} else {
			var str string
			k := 0
			t := 1
			var file *os.File
			for _, payment := range payments {
				if k == 0 {
					file, _ = os.OpenFile(dir+"/payments"+fmt.Sprint(t)+".dump", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
				}
				k++
				str = payment.ID + ";" + strconv.Itoa(int(payment.AccountID)) + ";" + strconv.Itoa(int(payment.Amount)) + ";" + string(payment.Category) + ";" + string(payment.Status) + "\n"
				_, err := file.WriteString(str)
				if err != nil {
					log.Print(err)
				}
				if k == records {
					str = ""
					t++
					k = 0
					fmt.Println(t, " = t")
					file.Close()
				}
			}
		}
	}
	return nil
}
func (s *Service) SumPayments(goroutines int) types.Money {
	money := types.Money(0)
	if goroutines < 2 {
		for _, value := range s.payments {
			money += value.Amount
		}
		return money
	}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	max := 0
	count := len(s.payments) / goroutines
	for i := 1; i < goroutines; i++ {
		max += count
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			sum := types.Money(0)
			for _, value := range s.payments[val-count : val] {
				sum += value.Amount
			}
			mutex.Lock()
			money += sum
			mutex.Unlock()
		}(max)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := types.Money(0)
		for _, value := range s.payments[max:] {
			sum += value.Amount
		}
		mutex.Lock()
		money += sum
		mutex.Unlock()
	}()
	wg.Wait()
	return money
	//wg := sync.WaitGroup{}
	//mu := sync.Mutex{}
	//sum := int64(0)
	//kol := 0
	//i := 0
	//if goroutines == 0 {
	//	kol = len(s.payments)
	//} else {
	//	kol = int(len(s.payments) / goroutines)
	//}
	//for i = 0; i < goroutines-1; i++ {
	//	wg.Add(1)
	//	go func(index int) {
	//		defer wg.Done()
	//		val := int64(0)
	//		payments := s.payments[index*kol : (index+1)*kol]
	//		for _, payment := range payments {
	//			val += int64(payment.Amount)
	//		}
	//		mu.Lock()
	//		sum += val
	//		mu.Unlock()
	//
	//	}(i)
	//}
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	val := int64(0)
	//	payments := s.payments[i*kol:]
	//	for _, payment := range payments {
	//		val += int64(payment.Amount)
	//	}
	//	mu.Lock()
	//	sum += val
	//	mu.Unlock()
	//
	//}()
	//wg.Wait()
	//return types.Money(sum)
}

func (s Service) FilterPayments(accountID int64, goroutines int) (newPayment []types.Payment, err error) {
	if goroutines < 2 {
		for _, value := range s.payments {
			if value.AccountID == accountID {
				newPayment = append(newPayment, *value)
			}
		}
		if newPayment == nil {
			return nil, ErrAccountNotFound
		}
		return
	}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	max := 0
	count := len(s.payments) / goroutines
	for i := 1; i < goroutines; i++ {
		max += count
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			sum := []types.Payment{}
			for _, value := range s.payments[val-count : val] {
				if value.AccountID == accountID {
					sum = append(sum, *value)
				}
			}
			mutex.Lock()
			newPayment = append(newPayment, sum...)
			mutex.Unlock()
		}(max)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := []types.Payment{}
		for _, value := range s.payments[max:] {
			if value.AccountID == accountID {
				sum = append(sum, *value)
			}
		}
		mutex.Lock()
		newPayment = append(newPayment, sum...)
		mutex.Unlock()
	}()
	wg.Wait()
	if newPayment == nil {
		return nil, ErrAccountNotFound
	}
	return
}

func (s Service) FilterPaymentsByFn(filter func(payment types.Payment) bool, goroutines int, ) (newPayment []types.Payment, err error) {
	if goroutines < 2 {
		for _, value := range s.payments {
			if filter(*value) {
				newPayment = append(newPayment, *value)
			}
		}
		if newPayment == nil {
			return nil, ErrAccountNotFound
		}
		return
	}
	wg := sync.WaitGroup{}
	mutex := sync.Mutex{}
	max := 0
	count := len(s.payments) / goroutines
	for i := 1; i < goroutines; i++ {
		max += count
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			sum := []types.Payment{}
			for _, value := range s.payments[val-count : val] {
				if filter(*value) {
					sum = append(sum, *value)
				}
			}
			mutex.Lock()
			newPayment = append(newPayment, sum...)
			mutex.Unlock()
		}(max)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		sum := []types.Payment{}
		for _, value := range s.payments[max:] {
			if filter(*value) {
				sum = append(sum, *value)
			}
		}
		mutex.Lock()
		newPayment = append(newPayment, sum...)
		mutex.Unlock()
	}()
	wg.Wait()
	if newPayment == nil {
		return nil, ErrAccountNotFound
	}
	return
}

func (s Service) SumPaymentsWithProgress1() <-chan Progress {
	data := make([]types.Payment, 1_000_000)
	for i := range data {
		data[i] = types.Payment{
			ID: "", AccountID: int64(i), Amount: types.Money(i), Category: "auto", Status: "ok",
		}
	}
	part := 10
	size := len(data) / part
	ch := make([]chan Progress, part)
	for i := 0; i < part; i++ {
		a := i
		c1 := make(chan Progress)
		ch[i] = c1
		go func(c chan<- Progress, d []types.Payment) {
			defer close(c)
			sum := Progress{}
			for _, j := range d {
				sum.Part = a
				sum.Result += j.Amount
			}
			c <- sum
		}(c1, data[i*size:(i+1)*size])
	}
	merge(ch)
	return nil
}
func merge(chanal []chan Progress) {
	wg := &sync.WaitGroup{}
	wg.Add(len(chanal))
	merge := make(chan Progress)
	for _, i := range chanal {
		go func(ch chan Progress) {
			defer wg.Done()
			for i := range ch {
				merge <- i
			}
		}(i)
	}
	go func() {
		defer close(merge)
		wg.Wait()
	}()
	total := types.Money(0)
	for i := range merge {
		total += i.Result
	}
}

func (s Service) SumPaymentsWithProgress() <-chan Progress {

	part := 10
	size := len(s.payments) / part
	wg := &sync.WaitGroup{}
	chanal := make(chan Progress, part)
	defer close(chanal)
	if s.payments == nil {
		return chanal
	}
	i := 0
	for i = 0; i < part; i++ {
		wg.Add(1)
		go func(ch chan Progress, j int) {
			defer wg.Done()
			sum := Progress{}
			for _, v := range s.payments[j*size : (j+1)*size] {
				sum.Result += v.Amount
			}
			ch <- sum
		}(chanal, i)
	}
	wg.Add(1)
	go func(ch chan Progress, j int) {
		defer wg.Done()
		sum := Progress{}
		for _, v := range s.payments[j*size : (j+1)*size] {
			sum.Result += v.Amount
		}
		ch <- sum
	}(chanal, i)

	wg.Wait()
	return chanal
}
