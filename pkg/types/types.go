package types

//Money ...
type Money int64

//PaymentCategory ...
type PaymentCategory string

//PaymentStatus ...
type PaymentStatus string

// ...
const (
	PaymentStatusOk         PaymentStatus = "Ok"
	PaymentStatusFail       PaymentStatus = "FAIL"
	PaymentStatusInProgress PaymentStatus = "INPROGRESS"
)
//Payment ...
type Payment struct {
	ID       string
	AccountID int64
	Amount   Money
	Category PaymentCategory
	Status   PaymentStatus
}

//Phone ...
type Phone string

//Account ...
type Account struct {
	ID      int64
	Phone   Phone
	Balance Money
}
type Favorite struct {
	ID string
	AccountID int64
	Name string
	Amount Money
	Category PaymentCategory
}