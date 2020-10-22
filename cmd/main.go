package main

import (
	"github.com/ilhom9045/wallet/pkg/types"
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
)

func main() {
	s := wallet.Service{}
	total := types.Money(0)

	for i := range s.SumPaymentsWithProgress1() {
		total += i.Result
	}
	
	log.Println(total)
}
