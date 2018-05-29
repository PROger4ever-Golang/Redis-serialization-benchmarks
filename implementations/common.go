package implementations

var KEY = "really_useful_object"

type SerializationObject struct {
	String1, String2, String3, String4, String5 string
	FieldX                                      string
}

type Implementation interface {
	Get() (object *SerializationObject, err error)
	Set(object *SerializationObject) (err error)
	Del() (err error)
	GetOneField(fieldName string) (value string, err error)
	SetOneField(fieldName string, value string) (err error)
}
