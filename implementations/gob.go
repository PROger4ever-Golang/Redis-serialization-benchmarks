package implementations

import (
	"bytes"
	"encoding/gob"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/gomodule/redigo/redis"
	"reflect"
)

type GobImplementation struct {
	connection redis.Conn
}

func (i *GobImplementation) Get() (object *SerializationObject, err error) {
	gobBytes, err := redis.Bytes(i.connection.Do("GET", KEY))
	if err != nil {
		return nil, utils.WrapIfError(err, "GobImplementation->Get()->GET")
	}

	object = new(SerializationObject)
	err = gob.NewDecoder(bytes.NewReader(gobBytes)).Decode(&object)
	return object, utils.WrapIfError(err, "GobImplementation->Get()->Decode()")
}

func (i *GobImplementation) Set(object *SerializationObject) (err error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(object)

	_, err = i.connection.Do("SET", KEY, buf.Bytes())
	return utils.WrapIfError(err, "GobImplementation->Set()->SET")
}

func (i *GobImplementation) Del() (err error) {
	_, err = i.connection.Do("DEL", KEY)
	return utils.WrapIfError(err, "GobImplementation->Del()->DEL")
}

func (i *GobImplementation) GetOneField(fieldName string) (value string, err error) {
	object, err := i.Get()
	if err != nil {
		return "", utils.WrapIfError(err, "GobImplementation->GetOneField()->Get()")
	}

	val := reflect.ValueOf(object).Elem()
	value = val.FieldByName(fieldName).String()
	return
}

func (i *GobImplementation) SetOneField(fieldName string, value string) (err error) {
	object, err := i.Get()
	if err != nil {
		return utils.WrapIfError(err, "GobImplementation->SetOneField()->Get()")
	}

	val := reflect.ValueOf(object).Elem()
	val.FieldByName(fieldName).SetString(value)
	return i.Set(object)
}

func NewGobImplementation(connection redis.Conn) (implementation *GobImplementation) {
	implementation = new(GobImplementation)
	implementation.connection = connection
	gob.Register(implementation)
	return
}
