package implementations

import (
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/gomodule/redigo/redis"
)

type FlatImplementation struct {
	connection redis.Conn
}

func (i *FlatImplementation) Get() (object *SerializationObject, err error) {
	structMap, err := redis.Values(i.connection.Do("HGETALL", KEY))
	if err != nil {
		return nil, utils.WrapIfError(err, "FlatImplementation->Get()->HGETALL")
	}

	object = new(SerializationObject)
	err = redis.ScanStruct(structMap, object)
	return object, utils.WrapIfError(err, "FlatImplementation->Get()->ScanStruct")
}

func (i *FlatImplementation) Set(object *SerializationObject) (err error) {
	_, err = i.connection.Do("HMSET", redis.Args{KEY}.AddFlat(object)...)
	return utils.WrapIfError(err, "FlatImplementation->Set()->HMSET")
}

func (i *FlatImplementation) Del() (err error) {
	_, err = i.connection.Do("DEL", KEY)
	return utils.WrapIfError(err, "FlatImplementation->Del()->DEL")
}

func (i *FlatImplementation) GetOneField(fieldName string) (value string, err error) {
	value, err = redis.String(i.connection.Do("HGET", KEY, fieldName))
	return value, utils.WrapIfError(err, "FlatImplementation->GetOneField()->HGET")
}

func (i *FlatImplementation) SetOneField(fieldName string, value string) (err error) {
	_, err = i.connection.Do("HSET", KEY, fieldName, value)
	return utils.WrapIfError(err, "FlatImplementation->SetOneField()->HSET")
}

func NewFlatImplementation(connection redis.Conn) (implementation *FlatImplementation) {
	implementation = new(FlatImplementation)
	implementation.connection = connection
	return
}
