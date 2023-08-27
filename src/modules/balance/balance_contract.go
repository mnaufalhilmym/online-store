package balance

type addBalanceReq struct {
	Amount *int `json:"amount" validate:"required,gt=0"`
}
