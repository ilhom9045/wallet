package wallet

import (
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
