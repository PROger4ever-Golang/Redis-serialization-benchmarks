package implementations

import (
	"encoding/json"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/gomodule/redigo/redis"
	"github.com/nitishm/go-rejson"
	"strings"
)

type RejsonImplementation struct {
	connection redis.Conn
}

func (i *RejsonImplementation) Get() (object *SerializationObject, err error) {
	jsonBytes, err := redis.Bytes(rejson.JSONGet(i.connection, KEY, ""))
	if err != nil {
		return nil, utils.WrapIfError(err, "RejsonImplementation->Get()->JSONGet")
	}

	object = new(SerializationObject)
	err = json.Unmarshal(jsonBytes, object)
	return object, utils.WrapIfError(err, "RejsonImplementation->Get()->Unmarshal()")
}

func (i *RejsonImplementation) Set(object *SerializationObject) (err error) {
	_, err = rejson.JSONSet(i.connection, KEY, ".", object, false, false)
	return utils.WrapIfError(err, "RejsonImplementation->Set()->JSONSet")
}

func (i *RejsonImplementation) Del() (err error) {
	_, err = i.connection.Do("DEL", KEY)
	return utils.WrapIfError(err, "RejsonImplementation->Del()->DEL")
}
func (i *RejsonImplementation) GetOneField(fieldName string) (value string, err error) {
	value, err = redis.String(rejson.JSONGet(i.connection, KEY, "."+fieldName))
	value = strings.Trim(value, `"`) //Note: X_x
	return value, utils.WrapIfError(err, "RejsonImplementation->GetOneField()->JSONGet")
}

func (i *RejsonImplementation) SetOneField(fieldName string, value string) (err error) {
	_, err = rejson.JSONSet(i.connection, KEY, "."+fieldName, value, false, false)
	return utils.WrapIfError(err, "RejsonImplementation->SetOneField()->JSONSet")
}

func NewRejsonImplementation(connection redis.Conn) (implementation *RejsonImplementation) {
	implementation = new(RejsonImplementation)
	implementation.connection = connection
	return
}
