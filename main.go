package main

import (
	"fmt"
	"github.com/PROger4ever-Golang/Redis-serialization-benchmarks/implementations"
	"reflect"
)

const ERROR_FORMAT = "%s: error occured while %s - %s"

func main() {
	app, err := NewApplication()
	if err != nil {
		app.ErrorLogger.Printf(ERROR_FORMAT, "main", "configuring app", err)
		return
	}
	defer app.Finalize()

	for name, imp := range app.Implementations {
		imp.Del()
		srcObject := &implementations.SerializationObject{
			String1: "String1", String2: "String2", String3: "String3", String4: "String4", String5: "String5",
			FieldX: "FieldX",
		}

		err = imp.Set(srcObject)
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "setting object", err)
			continue
		}

		object1, err := imp.Get()
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "getting object", err)
			continue
		}

		isEqual := reflect.DeepEqual(srcObject, object1)
		if !isEqual {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "comparing src srcObject and object1", "they aren't equal")
			continue
		}

		srcObject.FieldX = "FieldX-changed"
		err = imp.SetOneField("FieldX", srcObject.FieldX)
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "setting one field", err)
			continue
		}

		field, err := imp.GetOneField("FieldX")
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "getting one field", err)
			continue
		}

		if field != srcObject.FieldX {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "comparing src field and actual field", "they aren't equal")
			continue
		}

		object2, err := imp.Get()
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "getting object", err)
			continue
		}

		isEqual = reflect.DeepEqual(srcObject, object2)
		if !isEqual {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "comparing src srcObject and object2", "they aren't equal")
			continue
		}

		err = imp.Del()
		if err != nil {
			app.ErrorLogger.Printf(ERROR_FORMAT, name, "deleting object", err)
			continue
		}

		fmt.Printf("%s, %+v\n", field, object1)
	}
}
