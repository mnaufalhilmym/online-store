package account

type updateAccountReq struct {
	Name     *string `json:"name"`
	Username *string `json:"username"`
	Password *string `json:"password"`
}
