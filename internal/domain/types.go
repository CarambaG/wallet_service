package domain

type OperationType string

const (
	OperationDeposit  OperationType = "DEPOSIT"
	OperationWithdraw OperationType = "WITHDRAW"
)

func (o OperationType) Valid() bool {
	switch o {
	case OperationDeposit, OperationWithdraw:
		return true
	default:
		return false
	}
}
