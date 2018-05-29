package implementations

import (
	"encoding/json"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/utils"
	"github.com/gomodule/redigo/redis"
)

type JsonImplementation struct {
	connection redis.Conn
}

func (i *JsonImplementation) Get() (object *SerializationObject, err error) {
	jsonString, err := redis.String(i.connection.Do("GET", KEY))
	if err != nil {
		return nil, utils.WrapIfError(err, "JsonImplementation->Get()->GET")
	}

	object = new(SerializationObject)
	err = json.Unmarshal([]byte(jsonString), object)
	return object, utils.WrapIfError(err, "JsonImplementation->Get()->Unmarshal()")
}

func (i *JsonImplementation) Set(object *SerializationObject) (err error) {
	jsonBytes, err := json.Marshal(object)
	if err != nil {
		return utils.WrapIfError(err, "JsonImplementation->Set()->Marshal()")
	}
	jsonString := string(jsonBytes)
	_, err = i.connection.Do("SET", KEY, jsonString)
	return utils.WrapIfError(err, "JsonImplementation->Set()->SET")
}

func (i *JsonImplementation) Del() (err error) {
	_, err = i.connection.Do("DEL", KEY)
	return utils.WrapIfError(err, "JsonImplementation->Del()->DEL")
}

func (i *JsonImplementation) getMap() (theMap map[string]string, err error) {
	jsonString, err := redis.String(i.connection.Do("GET", KEY))
	if err != nil {
		return nil, utils.WrapIfError(err, "JsonImplementation->getMap()->GET")
	}

	theMap = make(map[string]string)
	err = json.Unmarshal([]byte(jsonString), &theMap)
	return theMap, utils.WrapIfError(err, "JsonImplementation->getMap()->Unmarshal()")
}

func (i *JsonImplementation) GetOneField(fieldName string) (value string, err error) {
	theMap, err := i.getMap()
	if err != nil {
		return "", utils.WrapIfError(err, "JsonImplementation->GetOneField()->getMap()")
	}

	if field, ok := theMap[fieldName]; ok {
		return field, nil
	}
	return "", utils.WrapIfError(err, "JsonImplementation->GetOneField()->theMap[fieldName]")
}

func (i *JsonImplementation) SetOneField(fieldName string, value string) (err error) {
	theMap, err := i.getMap()
	if err != nil {
		return utils.WrapIfError(err, "JsonImplementation->SetOneField->getMap()")
	}
	theMap[fieldName] = value

	jsonBytes, err := json.Marshal(theMap)
	if err != nil {
		return utils.WrapIfError(err, "JsonImplementation->SetOneField->Marshal()")
	}
	jsonString := string(jsonBytes)
	_, err = i.connection.Do("SET", KEY, jsonString)
	return utils.WrapIfError(err, "JsonImplementation->SetOneField->SET")
}

func NewJsonImplementation(connection redis.Conn) (implementation *JsonImplementation) {
	implementation = new(JsonImplementation)
	implementation.connection = connection
	return
}
