package main

import (
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

	svc := &wallet.Service{}

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
	i := svc.SumPayments(9)
	log.Println(i)

}
