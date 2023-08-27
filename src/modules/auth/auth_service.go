package auth

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	a "hilmy.dev/store/src/modules/account/account_entity"
	b "hilmy.dev/store/src/modules/balance/balance_entity"
)

func (m *Module) getAccountDetailService(id *uuid.UUID) (*a.AccountModel, error) {
	return a.AccountRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (m *Module) getAccountDetailByUsernameService(username *string) (*a.AccountModel, error) {
	return a.AccountRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "username = ?",
				Args:  []interface{}{username},
			},
		},
	})
}

func (m *Module) addAccountService(data *a.AccountModel) (*a.AccountModel, error) {
	return a.AccountRepository().Create(data)
}

func (m *Module) deleteAccountService(id *uuid.UUID) error {
	return a.AccountRepository().Destroy(&a.AccountModel{
		Model: pg.Model{
			ID: id,
		},
	})
}

func (m *Module) createBalanceService(data *b.BalanceModel) (*b.BalanceModel, error) {
	return b.BalanceRepository().Create(data)
}
