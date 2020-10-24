package wallet

import (
	"github.com/ilhom9045/wallet/pkg/types"
	"log"
	"testing"
)

func TestService_FindAccountByID_true(t *testing.T) {

	service := Service{}
	service.RegisterAccount("+992927459045")
	_, err := service.FindAccountByID(1)
	if err != nil {
		t.Error(err)
	}
}
func TestService_FindAccountByID_false(t *testing.T) {

	service := Service{}
	service.RegisterAccount("+992927459045")
	_, err := service.FindAccountByID(7)
	if err == nil {
		t.Error(err)
	}
}
func TestService_Reject_success(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Reject(pay.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Reject_fail(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
	err = svc.Reject(pay.ID + "asd")
	if err == nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}
}

func TestService_Repeat_true(t *testing.T) {
	svc := Service{}
	svc.RegisterAccount("+9920000001")

	account, err := svc.FindAccountByID(1)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	err = svc.Deposit(account.ID, 1000_00)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	payment, err := svc.Pay(account.ID, 100_00, "auto")
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err := svc.FindPaymentByID(payment.ID)
	if err != nil {
		t.Errorf("\ngot > %v \nwant > nil", err)
	}

	pay, err = svc.Repeat(pay.ID)
	if err != nil {
		t.Errorf("Repeat(): Error(): can't pay for an account(%v): %v", pay.ID, err)
	}
}
func TestService_Favorite_success_user(t *testing.T) {
	svc := Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
		t.Errorf("method RegisterAccount returned not nil error, account => %v", account)
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		t.Errorf("method Deposit returned not nil error, error => %v", err)
	}

	payment, err := svc.Pay(account.ID, 10_00, "auto")
	if err != nil {
		t.Errorf("Pay() Error() can't pay for an account(%v): %v", account, err)
	}

	favorite, err := svc.FavoritePayment(payment.ID, "megafon")
	if err != nil {
		t.Errorf("FavoritePayment() Error() can't for an favorite(%v): %v", favorite, err)
	}

	paymentFavorite, err := svc.PayFromFavorite(favorite.ID)
	if err != nil {
		t.Errorf("PayFromFavorite() Error() can't for an favorite(%v): %v", paymentFavorite, err)
	}
}
func TestService_SumPayments(b *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	want := types.Money(7000)

	got := svc.SumPayments(5)
	if want != got {
		b.Errorf(" error, want => %v got => %v", want, got)
	}

}
func TestService_ExportImport_success_user(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("+992000000001")
	err := svc.Deposit(account.ID, 100_00)
	if err != nil {
		log.Println(err, "13")
		return
	}

	account, _ = svc.RegisterAccount("+992000000002")
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		log.Println(err, 26)
	}

	dir, _ := svc.GetDir()
	//os.MkdirAll(dir,0777)
	err = svc.ExportToFile("data/export.txt")
	if err != nil {
		log.Println(err, 39)
	}
	err = svc.ImportFromFile("data/export.txt")
	if err != nil {
		log.Println(err, 43)
	}
	err = svc.Export(dir)
	if err != nil {
		log.Println(err, 47)
	}
	err = svc.Import(dir)
	if err != nil {
		log.Print(err, 51)
		t.Error(err)
	}
}
func TestService_Export_success_user(t *testing.T) {
	svc := &Service{}
	account, _ := svc.RegisterAccount("+992000000001")
	err := svc.Deposit(account.ID, 100_00)
	if err != nil {
		log.Println(err, "13")
		return
	}

	account, _ = svc.RegisterAccount("+992000000002")
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
		log.Println(err, 26)
	}

	dir, _ := svc.GetDir()
	//os.MkdirAll(dir,0777)
	err = svc.ExportToFile("data/export.txt")
	if err != nil {
		log.Println(err, 39)
	}
	err = svc.ImportFromFile("data/export.txt")
	if err != nil {
		log.Println(err, 43)
	}
	err = svc.Export(dir)
	if err != nil {
		log.Println(err, 47)
	}
	err = svc.Import(dir)
	if err != nil {
		log.Print(err, 51)
		return
	}

}
func TestService_ExportHistory_success_user(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	payment, err := svc.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = svc.HistoryToFiles(payment, "data", 4)
	if err != nil {
		t.Error(err)
	}
}
func TestService_ExportHistory(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000001")
	if err != nil {
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	payment, err := svc.ExportAccountHistory(1)
	if err != nil {
		t.Error(err)
	}
	err = svc.HistoryToFiles(payment, "data", 20)
	if err != nil {
		t.Error(err)
	}
}
func TestService_ExportToFile(t *testing.T) {
	s := Service{}
	err := s.ExportToFile("export.txt")
	if err != nil {
		t.Error(err)
	}
}
func TestService_ImportFromFile(t *testing.T) {
	s := Service{}
	err := s.ImportFromFile("export.txt")
	if err != nil {
		t.Error(err)
	}
}
func BenchmarkService_SumPayments(b *testing.B) {

	svc := Service{}

	account, err := svc.RegisterAccount("+992000000001")

	if err != nil {
	}

	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	_, err = svc.Pay(account.ID, 1, "Cafe")
	_, err = svc.Pay(account.ID, 2, "Cafe")
	_, err = svc.Pay(account.ID, 3, "Cafe")
	_, err = svc.Pay(account.ID, 4, "Cafe")
	_, err = svc.Pay(account.ID, 5, "Cafe")
	_, err = svc.Pay(account.ID, 6, "Cafe")
	_, err = svc.Pay(account.ID, 7, "Cafe")
	_, err = svc.Pay(account.ID, 8, "Cafe")
	_, err = svc.Pay(account.ID, 9, "Cafe")
	_, err = svc.Pay(account.ID, 10, "Cafe")
	_, err = svc.Pay(account.ID, 11, "Cafe")
	_, err = svc.Pay(account.ID, 12, "Cafe")
	_, err = svc.Pay(account.ID, 13, "Cafe")
	if err != nil {
	}
	for i := 0; i < b.N; i++ {
		svc.SumPayments(5)
	}

}
func BenchmarkService_FilterPayments(b *testing.B) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000000")
	account1, err := svc.RegisterAccount("+992000000001")
	account2, err := svc.RegisterAccount("+992000000002")
	account3, err := svc.RegisterAccount("+992000000003")
	account4, err := svc.RegisterAccount("+992000000004")
	acc, err := svc.RegisterAccount("+992000000005")
	if err != nil {
	}
	svc.Deposit(acc.ID, 100)
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")

	a, err := svc.FilterPayments(account.ID, 5)
	if err != nil {
		b.Error(err)
	}
	log.Println(len(a))
}
func BenchmarkService_FilterPaymentsByFn(b *testing.B) {
	svc := &Service{}
	filter := func(payment types.Payment) bool {
		for _, value := range svc.payments {
			if payment.ID == value.ID {
				return true
			}
		}
		return false
	}
	account, err := svc.RegisterAccount("+992000000000")
	account1, err := svc.RegisterAccount("+992000000001")
	account2, err := svc.RegisterAccount("+992000000002")
	account3, err := svc.RegisterAccount("+992000000003")
	account4, err := svc.RegisterAccount("+992000000004")
	acc, err := svc.RegisterAccount("+992000000005")
	if err != nil {
	}
	svc.Deposit(acc.ID, 100)
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	a, err := svc.FilterPaymentsByFn(filter, 4)
	if err != nil {
		b.Error(err)
	}
	log.Println(a)
}

//func BenchmarkService_FilterPaymentsByFn2(b *testing.B) {
//	svc := &Service{}
//	filter := func(payment types.Payment) bool {
//		for _, value := range svc.payments {
//			if payment.ID == value.ID {
//				return true
//			}
//		}
//		return false
//	}
//	account, err := svc.RegisterAccount("+992000000000")
//	account1, err := svc.RegisterAccount("+992000000001")
//	account2, err := svc.RegisterAccount("+992000000002")
//	account3, err := svc.RegisterAccount("+992000000003")
//	account4, err := svc.RegisterAccount("+992000000004")
//	acc, err := svc.RegisterAccount("+992000000005")
//	if err != nil {
//	}
//	svc.Deposit(acc.ID, 100)
//	err = svc.Deposit(account.ID, 100_00)
//	err = svc.Deposit(account1.ID, 100_00)
//	err = svc.Deposit(account2.ID, 100_00)
//	err = svc.Deposit(account3.ID, 100_00)
//	err = svc.Deposit(account4.ID, 100_00)
//	err = svc.Deposit(account2.ID, 100_00)
//	err = svc.Deposit(account3.ID, 100_00)
//	err = svc.Deposit(account4.ID, 100_00)
//	if err != nil {
//	}
//
//	_, err = svc.FilterPaymentsByFn(filter, 4)
//	if err == nil {
//		b.Error(err)
//	}
//}
func BenchmarkService_SumPaymentsWithProgress(b *testing.B) {
	s := Service{}
	for i := 0; i < 1_000_001; i++ {
		s.Pay(1, types.Money(i), "auto")
	}
	for i := 0; i < b.N; i++ {
		s.SumPaymentsWithProgress()
	}
}

//func BenchmarkService_SumPaymentsWithProgress1(b *testing.B) {
//	s := Service{}
//	for i := 0; i < b.N; i++ {
//		s.SumPaymentsWithProgress2()
//	}
//}
func TestService_FilterPayments(t *testing.T) {
	svc := &Service{}

	account, err := svc.RegisterAccount("+992000000000")
	account1, err := svc.RegisterAccount("+992000000001")
	account2, err := svc.RegisterAccount("+992000000002")
	account3, err := svc.RegisterAccount("+992000000003")
	account4, err := svc.RegisterAccount("+992000000004")
	acc, err := svc.RegisterAccount("+992000000005")
	if err != nil {
	}
	svc.Deposit(acc.ID, 100)
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")

	a, err := svc.FilterPayments(account.ID, 5)
	if err != nil {
		t.Error(err)
	}
	log.Println(len(a))
}
func TestService_FilterPaymentsByFn(t *testing.T) {
	svc := &Service{}
	filter := func(payment types.Payment) bool {
		for _, value := range svc.payments {
			if payment.ID == value.ID {
				return true
			}
		}
		return false
	}
	account, err := svc.RegisterAccount("+992000000000")
	account1, err := svc.RegisterAccount("+992000000001")
	account2, err := svc.RegisterAccount("+992000000002")
	account3, err := svc.RegisterAccount("+992000000003")
	account4, err := svc.RegisterAccount("+992000000004")
	acc, err := svc.RegisterAccount("+992000000005")
	if err != nil {
	}
	svc.Deposit(acc.ID, 100)
	err = svc.Deposit(account.ID, 100_00)
	if err != nil {
	}

	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account1.ID, 10_00, "auto")
	svc.Pay(account3.ID, 10_00, "auto")
	svc.Pay(account2.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	svc.Pay(account4.ID, 10_00, "auto")
	a, err := svc.FilterPaymentsByFn(filter, 4)
	if err != nil {
		t.Error(err)
	}
	log.Println(a)
}

