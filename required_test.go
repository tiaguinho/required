package required_test

import (
	"testing"
	"github.com/tiaguinho/required"
	"reflect"
)

var (
	c = required.Message{
		Field: "C",
		ErrMsg: "where is the value?",
	}

	d = required.Message{
		Field: "default",
		ErrMsg: "this field is required",
	}
)

type I struct {
	C    int  `required:"where is the number?"`
	D    int  `json:"default" required:"-"`
	N    int  `json:"do_not_check"`
}

func TestValidate(t *testing.T) {
	v := I{}

	if err := required.Validate(v); err != nil {
		sm := make([]required.Message, 2)
		sm[0] = c
		sm[1] = d

		e := required.New(sm...)
		if e == err {
			t.Errorf("\n expected: \n %s \n got: \n %s", e, err)
		}
	}

	v.D = 100
	v.C = 100

	err := required.Validate(v)
	if err != nil {
		t.Errorf("\n no error message expected \n got: \n %s", err)
	}
}

func TestValidateWithMessage(t *testing.T) {
	v := I{}

	if err, msg := required.ValidateWithMessage(v); err != nil {
		sm := make([]required.Message, 0)
		sm = append(sm, d, c)

		if reflect.DeepEqual(sm, msg) {
			t.Errorf("\n expected: \n %s \n got: \n %s", sm, msg)
		}
	}

	v.D = 100
	v.C = 100

	err, msg := required.ValidateWithMessage(v)
	if err != nil || len(msg) > 0 {
		t.Errorf("\n no error message expected \n got: \n %s \n %+v", err, msg)
	}
}