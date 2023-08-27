package balance

import (
	"github.com/google/uuid"
	"hilmy.dev/store/src/libs/db/pg"
	b "hilmy.dev/store/src/modules/balance/balance_entity"
)

func (*Module) getBalanceByUserIDService(userID *uuid.UUID) (*b.BalanceModel, error) {
	return b.BalanceRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ?",
				Args:  []interface{}{userID},
			},
		},
	})
}

func (m *Module) updateBalanceByUserIDService(userID *uuid.UUID, data *b.BalanceModel) (*b.BalanceModel, error) {
	if _, err := b.BalanceRepository().Update(data, &pg.UpdateOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ?",
				Args:  []interface{}{userID},
			},
		},
	}); err != nil {
		return nil, err
	}

	data, err := b.BalanceRepository().FindOne(&pg.FindOneOptions{
		Where: &[]pg.Where{
			{
				Query: "user_id = ?",
				Args:  []interface{}{userID},
			},
		},
	})
	if err != nil {
		return nil, err
	}

	return data, nil
}
