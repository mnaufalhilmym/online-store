package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"hilmy.dev/store/src/libs/validator"
)

type ModelI interface {
	DatabaseName() string
	CollectionName() string
}

type Service[T ModelI] struct {
	Client *Client
}

type Options struct {
	UniqueFields []string
	Expiration   time.Duration
}

func NewService[T ModelI](client *Client, moptions ...*Options) *Service[T] {
	if client == nil {
		logger.Panic("db cannot be nil")
	}

	model := new(T)

	if len(moptions) > 0 {
		indexModels := func() []mongo.IndexModel {
			model := make([]mongo.IndexModel, 0)

			if moptions[0].Expiration > 0 {
				model = append(model, mongo.IndexModel{
					Keys:    bson.M{"updated_at": 1},
					Options: options.Index().SetExpireAfterSeconds(int32(moptions[0].Expiration / time.Second)),
				})
			}

			if len(moptions[0].UniqueFields) > 0 {
				uniqueField := bson.D{}
				for _, field := range moptions[0].UniqueFields {
					uniqueField = append(uniqueField, bson.E{Key: field, Value: 1})
				}
				model = append(model, mongo.IndexModel{
					Keys:    uniqueField,
					Options: options.Index().SetUnique(true),
				})
			}

			return model
		}()

		client.Database((*model).DatabaseName()).Collection((*model).CollectionName()).Indexes().CreateMany(context.TODO(), indexModels)
	}

	return &Service[T]{Client: client}
}

func (s *Service[T]) Count(countOptions *CountOptions) (*int64, error) {
	model := new(T)

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	count, err := coll.CountDocuments(context.TODO(), countOptions.Where)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return &count, nil
}

func (s *Service[T]) FindOne(findOptions *FindOneOptions) (*T, error) {
	model := new(T)

	docStruct := new(T)

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	docMap := make(map[string]interface{}, 0)
	if err := coll.FindOne(context.TODO(), findOptions.Where).Decode(&docMap); err != nil {
		if !IsErrNoDocuments(err) {
			logger.Error(err)
		}
		return nil, err
	}
	if err := transformInterfaceToStruct(&docMap, &docStruct); err != nil {
		return nil, err
	}

	return docStruct, nil
}

func (s *Service[T]) FindAll(findOptions *FindAllOptions) (*[]*T, *Pagination, error) {
	model := new(T)

	docStruct := []*T{}

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())
	optsFind := options.Find()

	if findOptions.Limit != nil && *findOptions.Limit > 0 {
		if *findOptions.Limit < FindAllMaximumLimit {
			optsFind = optsFind.SetLimit(int64(*findOptions.Limit))
		} else {
			optsFind = optsFind.SetLimit(FindAllMaximumLimit)
		}
	} else {
		optsFind = optsFind.SetLimit(FindAllDefaultLimit)
	}

	where := []Where{}
	if findOptions.Where != nil {
		for i := range *findOptions.Where {
			if (*findOptions.Where)[i].IncludeInCount {
				where = append(where, (*findOptions.Where)[i].Where)
			}
		}
	}

	total, err := coll.CountDocuments(context.TODO(), where)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	if findOptions.Offset != nil && *findOptions.Offset > 0 {
		optsFind = optsFind.SetSkip(int64(*findOptions.Offset))
	} else {
		optsFind = optsFind.SetSkip(0)
	}

	if findOptions.Where != nil {
		for i := range *findOptions.Where {
			if !(*findOptions.Where)[i].IncludeInCount {
				where = append(where, (*findOptions.Where)[i].Where)
			}
		}
	}

	if findOptions.Order != nil {
		order := bson.D{}
		for i := range *findOptions.Order {
			order = append(order, bson.E{
				Key: (*findOptions.Order)[i].Key, Value: (*findOptions.Order)[i].Value,
			})
		}
		optsFind = optsFind.SetSort(&order)
	}

	cursor, err := coll.Find(context.TODO(), where, optsFind)
	if err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	docMapList := []map[string]interface{}{}
	for cursor.Next(context.TODO()) {
		docMap := make(map[string]interface{}, 0)
		if err := cursor.Decode(&docMap); err != nil {
			logger.Error(err)
			return nil, nil, err
		}
		docMapList = append(docMapList, docMap)
	}
	if len(docMapList) > 0 {
		if err := transformInterfaceToStruct(&docMapList, &docStruct); err != nil {
			return nil, nil, err
		}
	}

	return &docStruct, &Pagination{
		Limit: int(*optsFind.Limit),
		Count: len(docStruct),
		Total: int(total),
	}, nil
}

func (s *Service[T]) Create(data *T) (*primitive.ObjectID, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	model := new(T)

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	docMap := make(map[string]interface{}, 0)
	if err := transformStructToMap(data, &docMap); err != nil {
		return nil, err
	}

	result, err := coll.InsertOne(context.TODO(), docMap)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	id := result.InsertedID.(primitive.ObjectID)

	return &id, nil
}

func (s *Service[T]) Update(data *T, updateOptions *UpdateOptions) error {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return err
	}

	model := new(T)

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	docMap := make(map[string]interface{}, 0)
	if err := transformStructToMap(data, &docMap); err != nil {
		return err
	}

	result, err := coll.UpdateOne(context.TODO(), updateOptions.Where, bson.D{{Key: "$set", Value: docMap}})
	if err != nil {
		logger.Error(err)
		return err
	}
	if result == nil || result.ModifiedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

func (s *Service[T]) Destroy(destroyOptions *DestroyOptions) error {
	model := new(T)

	coll := s.Client.Database((*model).DatabaseName()).Collection((*model).CollectionName())

	result, err := coll.DeleteOne(context.TODO(), destroyOptions.Where)
	if err != nil {
		logger.Error(err)
		return err
	}
	if result == nil || result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
