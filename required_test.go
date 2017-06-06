package required_test

import (
	"fmt"
	"github.com/tiaguinho/required"
	"reflect"
	"testing"
)

var (
	i = required.Message{
		Field:  "I",
		ErrMsg: "where is the value?",
	}

	s = required.Message{
		Field:  "default",
		ErrMsg: "this field is required",
	}
)

type T struct {
	I int                    `required:"where is the number?"`
	U uint8                  `xml:"uint" required:"not empty"`
	F float32                `xml:"float" required:"-"`
	C complex64              `xml:"complex" required:"-"`
	S string                 `xml:"default" required:"-"`
	M map[string]interface{} `json:"map" required:"-"`
	A []A                    `json:"array" required:"-"`
	P *A                     `json:"ptr" required:"pointer is empty"`
	N int                    `json:"do_not_check"`
}

type A struct {
	I int    `required:"-"`
	S string `json:"s" required:"don't left this field blank!'"`
}

func TestError(t *testing.T) {
	e := required.New(required.Message{Index: 0, Field: "test", ErrMsg: "test message"})
	if e.Error() != fmt.Sprintf("[%s]: %s \n", "test", "test message") {
		t.Errorf("got: %s", e)
	}

}

func TestValidate(t *testing.T) {
	err := required.Validate("")
	if err == nil {
		t.Errorf("error expected! returned nil.%s", err)
	}

	v := T{}

	err = required.Validate(v)
	if err == nil {
		t.Error("error expected! returned nil.")
	}

	sm := make([]required.Message, 2)
	sm[0] = i
	sm[1] = s

	e := required.New(sm...)
	if e == err {
		t.Errorf("\n expected: \n %s \n got: \n %s", e, err)
	}

	v.A = make([]A, 1)
	v.A[0] = A{
		I: 50,
	}

	err = required.Validate(v)
	if err == nil {
		t.Errorf("error expected! returned nil.%s", err)
	}

	v.I = 100
	v.U = 1
	v.F = 10.5
	v.C = 1
	v.S = "ok"
	v.M = map[string]interface{}{"test": "ok"}
	v.A[0].S = "sub message"
	v.P = &A{I: 10, S: "TESTE"}

	err = required.Validate(v)
	if err != nil {
		t.Errorf("\n no error message expected \n got: \n %s", err)
	}
}

func TestValidateWithMessage(t *testing.T) {
	msg, err := required.ValidateWithMessage("")
	if err == nil {
		t.Errorf("error expected! returned nil.%s", err)
	}

	v := T{}

	msg, err = required.ValidateWithMessage(v)
	if err != nil {
		sm := make([]required.Message, 2)
		sm[0] = i
		sm[1] = s

		if reflect.DeepEqual(sm, msg) {
			t.Errorf("\n expected: \n %s \n got: \n %s", sm, msg)
		}
	}

	v.A = make([]A, 1)
	v.A[0] = A{
		I: 50,
	}

	err = required.Validate(v)
	if err == nil {
		t.Error("error expected! returned nil.")
	}

	v.I = 100
	v.U = 1
	v.F = 10.5
	v.C = 1
	v.S = "ok"
	v.M = map[string]interface{}{"test": "ok"}
	v.A[0].S = "sub message"
	v.P = &A{I: 10, S: "TESTE"}

	msg, err = required.ValidateWithMessage(v)
	if err != nil || len(msg) > 0 {
		t.Errorf("\n no error message expected \n got: \n %s \n %+v", err, msg)
	}
}
