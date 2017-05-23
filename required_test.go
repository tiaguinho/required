package required

import (
	"testing"
	"reflect"
)

var (
	defaultMessage = "this field is required"
)

type I struct {
	Default    int  `json:"default" required:"-"`
	Custom     int  `required:"where is the number?"`
}

//
func TestValidateInt(t *testing.T) {
	v := I{}

	dm := Message {
		Field: "default",
		Message: defaultMessage,
	}

	cm := Message{
		Field: "Custom",
		Message: "where is the number?",
	}

	if msgs, err := Validate(v); err != nil {
		arr := make([]Message, 0)
		arr = append(arr, dm, cm)

		if !reflect.DeepEqual(arr, msgs) {
			t.Errorf("\n expected: %+v \n got: %+v", arr, msgs)
		}
	}

	v.Default = 100
	v.Custom = 0
	if msgs, err := Validate(v); err != nil {
		arr := make([]Message, 0)
		arr = append(arr, cm)

		if !reflect.DeepEqual(arr, msgs) {
			t.Errorf("\n expected: %+v \n got: %+v", arr, msgs)
		}
	}

	v.Default = 100
	v.Custom = 100
	if msgs, err := Validate(v); err != nil {
		t.Errorf("\n no error message expected \n got: %+v", msgs)
	}
}
