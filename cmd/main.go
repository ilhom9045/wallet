package main

import (
	"github.com/ilhom9045/wallet/pkg/types"
	"github.com/ilhom9045/wallet/pkg/wallet"
	"log"
)

func main() {
	s := wallet.Service{}
	total := types.Money(0)
	for i := range s.SumPaymentsWithProgress() {
		total += i.Result
	}
	//499999500000
	//1000001
	//1000001
	log.Println(total)

}
