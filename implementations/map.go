package implementations

import (
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/fatih/structs"
	"github.com/gomodule/redigo/redis"
	"github.com/mitchellh/mapstructure"
)

type MapImplementation struct {
	connection redis.Conn
}

func (i *MapImplementation) Get() (object *SerializationObject, err error) {
	structMap, err := redis.StringMap(i.connection.Do("HGETALL", KEY))
	if err != nil {
		return nil, utils.WrapIfError(err, "MapImplementation->Get()->HGETALL")
	}

	object = new(SerializationObject)
	err = mapstructure.Decode(structMap, object)
	return object, utils.WrapIfError(err, "MapImplementation->Get()->Decode()")
}

func (i *MapImplementation) Set(object *SerializationObject) (err error) {
	structMap := structs.Map(object)
	args := make([]interface{}, 0, 1+len(structMap)*2)

	args = append(args, KEY)
	for k, v := range structMap {
		args = append(args, k, v)
	}
	_, err = i.connection.Do("HMSET", args...)
	return utils.WrapIfError(err, "MapImplementation->Set()->HMSET")
}

func (i *MapImplementation) Del() (err error) {
	_, err = i.connection.Do("DEL", KEY)
	return utils.WrapIfError(err, "MapImplementation->Del()->DEL")
}

func (i *MapImplementation) GetOneField(fieldName string) (value string, err error) {
	value, err = redis.String(i.connection.Do("HGET", KEY, fieldName))
	return value, utils.WrapIfError(err, "MapImplementation->GetOneField()->HGET")
}

func (i *MapImplementation) SetOneField(fieldName string, value string) (err error) {
	_, err = i.connection.Do("HSET", KEY, fieldName, value)
	return utils.WrapIfError(err, "MapImplementation->SetOneField()->HSET")
}

func NewMapImplementation(connection redis.Conn) (implementation *MapImplementation) {
	implementation = new(MapImplementation)
	implementation.connection = connection
	return
}
