package pg

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hilmy.dev/store/src/libs/validator"
)

type ModelI interface {
	TableName() string
}

type Service[T ModelI] struct {
	DB *DB
}

func NewService[T ModelI](db *DB) *Service[T] {
	if db == nil {
		logger.Panic("service.DB must exist")
	}

	model := new(T)
	if err := db.AutoMigrate(model); err != nil {
		logger.Panic(err)
	}

	return &Service[T]{DB: db}
}

func Transaction(db *DB, txs ...func(tx *DB) *DB) error {
	if err := db.Transaction(func(tx *DB) error {
		for i := range txs {
			if err := txs[i](tx).Error; err != nil {
				logger.Error(err)
				return err
			}
		}
		return nil
	}); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *Service[T]) Count(countOptions *CountOptions) (*int64, error) {
	docStruct := new(T)

	countQuery := s.DB.Model(docStruct)

	if countOptions.Where != nil {
		for _, where := range *countOptions.Where {
			countQuery = countQuery.Where(where.Query, where.Args...)
		}
	}
	if countOptions.IsUnscoped {
		countQuery = countQuery.Unscoped()
	}

	count := new(int64)
	if err := countQuery.Count(count).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return count, nil
}

func (s *Service[T]) FindOne(findOptions *FindOneOptions) (*T, error) {
	docStruct := new(T)

	selectQuery := s.DB.Model(docStruct)

	if findOptions.IncludeTables != nil {
		for _, table := range *findOptions.IncludeTables {
			selectQuery = selectQuery.Preload(table.Query, table.Args...)
		}
	}
	if findOptions.Where != nil {
		for _, where := range *findOptions.Where {
			selectQuery = selectQuery.Where(where.Query, where.Args...)
		}
	}
	if findOptions.Order != nil {
		for _, order := range *findOptions.Order {
			selectQuery = selectQuery.Order(order)
		}
	}
	if findOptions.IsUnscoped {
		selectQuery = selectQuery.Unscoped()
	}

	if err := selectQuery.Take(docStruct).Error; err != nil {
		if !IsErrRecordNotFound(err) {
			logger.Error(err)
		}
		return nil, err
	}

	return docStruct, nil
}

func (s *Service[T]) FindAll(findOptions *FindAllOptions) (*[]*T, *Pagination, error) {
	docStruct := &[]*T{}

	selectQuery := s.DB.Model(docStruct)

	if findOptions.IncludeTables != nil {
		for _, table := range *findOptions.IncludeTables {
			selectQuery = selectQuery.Preload(table.Query, table.Args...)
		}
	}
	if findOptions.Where != nil {
		for _, where := range *findOptions.Where {
			if where.IncludeInCount {
				selectQuery = selectQuery.Where(where.Where.Query, where.Where.Args...)
			}
		}
	}
	if findOptions.IsUnscoped {
		selectQuery = selectQuery.Unscoped()
	}

	var total int64
	selectQuery.Count(&total)

	if findOptions.Where != nil {
		for _, where := range *findOptions.Where {
			if !where.IncludeInCount {
				selectQuery = selectQuery.Where(where.Where.Query, where.Where.Args...)
			}
		}
	}
	if findOptions.Order != nil {
		for _, order := range *findOptions.Order {
			selectQuery = selectQuery.Order(order)
		}
	}
	if findOptions.Limit != nil && *findOptions.Limit > 0 {
		if *findOptions.Limit > FindAllMaximumLimit {
			*findOptions.Limit = FindAllMaximumLimit
		}
	} else {
		*findOptions.Limit = FindAllDefaultLimit
	}
	selectQuery = selectQuery.Limit(*findOptions.Limit)

	if findOptions.Offset != nil && *findOptions.Offset > 0 {
		selectQuery = selectQuery.Offset(*findOptions.Offset)
	}

	if err := selectQuery.Find(docStruct).Error; err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return docStruct, &Pagination{
		Limit: *findOptions.Limit,
		Count: len(*docStruct),
		Total: int(total),
	}, nil
}

func (s *Service[T]) Create(data *T, createOptions ...*CreateOptions) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := s.CreateTx(s.DB, data, createOptions...).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) BulkCreate(data *[]*T, createOptions ...*CreateOptions) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	if err := s.BulkCreateTx(s.DB, data, createOptions...).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) Update(data *T, updateOptions ...*UpdateOptions) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	tx := s.UpdateTx(s.DB, data, updateOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (s *Service[T]) BulkUpdate(data *[]*T, updateOptions ...*UpdateOptions) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	tx := s.BulkUpdateTx(s.DB, data, updateOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (s *Service[T]) Replace(data *T, replaceOptions ...*ReplaceOptions) error {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return err
	}

	tx := s.ReplaceTx(s.DB, data, replaceOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) Destroy(data *T, destroyOptions ...*DestroyOptions) error {
	tx := s.DestroyTx(s.DB, data, destroyOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) BulkDestroy(data *[]*T, destroyOptions ...*DestroyOptions) error {
	tx := s.BulkDestroyTx(s.DB, data, destroyOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) CreateTx(tx *DB, data *T, createOptions ...*CreateOptions) *DB {
	docStruct := new(T)

	insertQuery := tx.Model(docStruct)

	if len(createOptions) > 0 {
		if createOptions[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	return insertQuery.Create(data)
}

func (s *Service[T]) BulkCreateTx(tx *DB, data *[]*T, createOptions ...*CreateOptions) *DB {
	docStruct := new(T)

	insertQuery := tx.Model(docStruct)

	if len(createOptions) > 0 {
		if createOptions[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	return insertQuery.Create(data)
}

func (s *Service[T]) UpdateTx(tx *DB, data *T, updateOptions ...*UpdateOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(updateOptions) > 0 && updateOptions[0].Where != nil {
		for _, where := range *updateOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if updateOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) BulkUpdateTx(tx *DB, data *[]*T, updateOptions ...*UpdateOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(updateOptions) > 0 && updateOptions[0].Where != nil {
		for _, where := range *updateOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if updateOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) ReplaceTx(tx *DB, data *T, replaceOptions ...*ReplaceOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(replaceOptions) > 0 && replaceOptions[0].Where != nil {
		for _, where := range *replaceOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if replaceOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) DestroyTx(tx *DB, data *T, destroyOptions ...*DestroyOptions) *DB {
	deleteQuery := tx

	if len(destroyOptions) > 0 && destroyOptions[0].Where != nil {
		for _, where := range *destroyOptions[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
		if destroyOptions[0].IsUnscoped {
			deleteQuery = deleteQuery.Unscoped()
		}
	}

	return deleteQuery.Delete(data)
}

func (s *Service[T]) BulkDestroyTx(tx *DB, data *[]*T, destroyOptions ...*DestroyOptions) *DB {
	deleteQuery := tx

	if len(destroyOptions) > 0 && destroyOptions[0].Where != nil {
		for _, where := range *destroyOptions[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
		if destroyOptions[0].IsUnscoped {
			deleteQuery = deleteQuery.Unscoped()
		}
	}

	return deleteQuery.Delete(data)
}
