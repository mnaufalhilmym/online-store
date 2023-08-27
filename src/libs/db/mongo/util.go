package mongo

import (
	"time"

	"github.com/bytedance/sonic"
	"go.mongodb.org/mongo-driver/mongo"
)

func IsErrNoDocuments(err error) bool {
	return err == mongo.ErrNoDocuments
}

func transformStructToMap[T any](data *T, target *map[string]interface{}) error {
	docBytes, err := sonic.ConfigFastest.Marshal(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := sonic.ConfigFastest.Unmarshal(docBytes, target); err != nil {
		logger.Error(err)
		return err
	}
	now := time.Now()
	(*target)["created_at"] = now
	(*target)["updated_at"] = now

	return nil
}

func transformInterfaceToStruct[T any](data interface{}, target *T) error {
	docBytes, err := sonic.ConfigFastest.Marshal(data)
	if err != nil {
		logger.Error(err)
		return err
	}

	if err := sonic.ConfigFastest.Unmarshal(docBytes, target); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
