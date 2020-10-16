package main

import (
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
)

func main() {
	svc := &wallet.Service{}
	account, _ := svc.RegisterAccount("+992000000001")
	svc.Deposit(account.ID, 100_00)
	payment, _ := svc.Pay(account.ID, 10_00, "auto")
	favorite, _ := svc.FavoritePayment(payment.ID, "megafon")
	svc.PayFromFavorite(favorite.ID)
	account, _ = svc.RegisterAccount("+992000000002")
	svc.Deposit(account.ID, 100_00)
	payment, _ = svc.Pay(account.ID, 10_00, "auto")
	favorite, _ = svc.FavoritePayment(payment.ID, "megafon")
	svc.PayFromFavorite(favorite.ID)
	dir, _ := svc.GetDir()
	//os.MkdirAll(dir,0777)
	svc.ExportToFile("data/export.txt")
	svc.ImportFromFile("data/export.txt")
	svc.Export(dir)
	err := svc.Import(dir)
	if err != nil {
		log.Print(err)
	}

	//svc := &wallet.Service{}
	//accountTest , err := svc.RegisterAccount("+992000000001")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//err = svc.Deposit(accountTest.ID, 100_000_00)
	//if err != nil {
	//	switch err {
	//	case wallet.ErrAmountMustBePositive:
	//		fmt.Println("Сумма должна быть положительной")
	//	case wallet.ErrAccountNotFound:
	//		fmt.Println("Аккаунт пользователя не найден")
	//	}
	//	return
	//}
	//fmt.Println(accountTest.Balance)
	//
	//err = svc.Deposit(accountTest.ID, 200_000_00)
	//if err != nil {
	//	switch err {
	//	case wallet.ErrAmountMustBePositive:
	//		fmt.Println("Сумма должна быть положительной")
	//	case wallet.ErrAccountNotFound:
	//		fmt.Println("Аккаунт пользователя не найден")
	//	}
	//	return
	//}
	//fmt.Println(accountTest.Balance)
	//
	//
	//newPay, err := svc.Pay(accountTest.ID,10_000_00,"auto")
	//newPay, err = svc.Pay(accountTest.ID,10_000_00,"food")
	//newPay, err = svc.Pay(accountTest.ID,10_000_00,"animal")
	//newPay, err = svc.Pay(accountTest.ID,10_000_00,"car")
	//newPay, err = svc.Pay(accountTest.ID,10_000_00,"restaurent")
	//fmt.Println(accountTest.Balance)
	//fmt.Println(newPay)
	//fmt.Println(err)
	//
	//fav, errFav := svc.FavoritePayment(newPay.ID, "Babilon")
	//fmt.Println(errFav)
	//fmt.Println(fav)
	//
	//wd, err := os.Getwd()
	//if err != nil {
	//	log.Print(err)
	//	return
	//}
	//err = svc.Export(wd)
	//if err != nil {
	//	return
	//}
	//err = svc.Import(wd)
	//if err != nil {
	// 	log.Print(err)
	// 	return
	//}

}
