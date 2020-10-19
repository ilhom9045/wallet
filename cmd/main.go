package main

import (
	"github.com/ilhom9045/wallet/pkg/types"
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
)

func main() {
	//svc := &wallet.Service{}
	//account, _ := svc.RegisterAccount("+992000000001")
	//err := svc.Deposit(account.ID, 100_00)
	//if err != nil {
	//	log.Println(err , "13")
	//	return
	//}
	//
	//account, _ = svc.RegisterAccount("+992000000002")
	//err = svc.Deposit(account.ID, 100_00)
	//if err != nil {
	//	log.Println(err,26)
	//}
	//
	//dir, _ := svc.GetDir()
	////os.MkdirAll(dir,0777)
	//err = svc.ExportToFile("data/export.txt")
	//if err != nil {
	//	log.Println(err,39)
	//}
	//err = svc.ImportFromFile("data/export.txt")
	//if err != nil {
	//	log.Println(err,43)
	//}
	//err = svc.Export(dir)
	//if err != nil {
	//	log.Println(err,47)
	//}
	//err = svc.Import(dir)
	//if err != nil {
	//	log.Print(err,51)
	//	return
	//}

	//svc := &wallet.Service{}
	//
	//account, err := svc.RegisterAccount("+992000000001")
	//if err != nil {
	//}
	//
	//err = svc.Deposit(account.ID, 100_00)
	//if err != nil {
	//}
	//
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//svc.Pay(account.ID, 10_00, "auto")
	//i := svc.SumPayments(9)
	//log.Println(i)
	svc := wallet.Service{}

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
	if err != nil {
	}

	want := types.Money(78)

	got := svc.SumPayments(2)
	log.Println(got)
	if want != got {
		log.Println("error sumpayment method")
	}

}
