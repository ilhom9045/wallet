package main

import (
	"github.com/ilhom9045/wallet/pkg/types"
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
)

func main() {
	s := wallet.Service{}
	total := types.Money(0)
	account, err := s.RegisterAccount("+992000000001")

	if err != nil {
	}

	err = s.Deposit(account.ID, 100_00)
	if err != nil {
	}
	for i := 0; i < 1_000_001; i++ {
		s.Pay(account.ID, types.Money(i), "auto")
	}

	log.Println(s.SumPayments(4))

	for i := range s.SumPaymentsWithProgress() {
		total += i.Result
	}
	//499999500000
	//1000001
	//1000001
	log.Println(total)

}
