package account

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	a "hilmy.dev/store/src/modules/account/account_entity"
)

func (*Module) getAccountDetailService(id *uuid.UUID) (*a.AccountModel, error) {
	return a.AccountRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
}

func (*Module) updateAccountService(id *uuid.UUID, data *a.AccountModel) (*a.AccountModel, error) {
	if _, err := a.AccountRepository().Update(data, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := a.AccountRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "id = ?",
				Args:  []interface{}{id},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (*Module) deleteAccountService(id *uuid.UUID) error {
	return a.AccountRepository().Destroy(&a.AccountModel{
		Model: pg.Model{
			ID: id,
		},
	})
}
