package main

import (
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
	"os"
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
	dir += "/export/"
	os.MkdirAll(dir,0777)
	svc.Export(dir)
	err := svc.Import(dir)
	if err != nil {
		log.Print(err)
	}
}
