package main

import (
	"github.com/ilhom9045/wallet/pkg/wallet"
)

func main() {
	service := &wallet.Service{}
	service.RegisterAccount("992927459045")
	service.Deposit(1, 1000)
	service.RegisterAccount("992927459046")
	// service.Deposit(1,1001)
	service.ExportToFile("data/export.txt")
	// service.ImportFromFile("data/import.txt")
}
