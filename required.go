package required

import (
	"fmt"
	"reflect"
)

// New returns an error with the messages passed in the param `sm`
func New(sm ...Message) error {
	return &errorMessages{sm}
}

// Error append all the fields errors
type errorMessages struct {
	sm []Message
}

// implement error interface
func (e *errorMessages) Error() (msg string) {
	for _, m := range e.sm {
		msg += fmt.Sprintf("[%s]: %s \n", m.Field, m.ErrMsg)
	}

	return
}

// Message hold the field name and message
// this can be return in some API implementation, where frontend can
// use this information to highlight the field and display the error message
type Message struct {
	Index  int    `json:"index,omitempty"`
	Field  string `json:"field" xml:"field"`
	ErrMsg string `json:"message" xml:"message"`
}

// Validate return error if any field is left empty
func Validate(v interface{}) error {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return fmt.Errorf("only struct can be validated")
	}

	sm := structFields(reflect.ValueOf(v))
	if len(sm) == 0 {
		return nil
	}
	return &errorMessages{sm}
}

// ValidateWithMessage return error and slice of message if any field is left empty
func ValidateWithMessage(v interface{}) ([]Message, error) {
	if reflect.TypeOf(v).Kind() != reflect.Struct {
		return []Message{}, fmt.Errorf("only struct can be validated")
	}

	sm := structFields(reflect.ValueOf(v))
	if len(sm) == 0 {
		return sm, nil
	}
	return sm, &errorMessages{sm}
}

// structFields check the type of field and if the field has a required tag
func structFields(v reflect.Value) []Message {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		if v.Kind() != reflect.Struct {
			return []Message{}
		}
	}

	m := make([]Message, 0)
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)

		if f.Kind() == reflect.Ptr {
			f = f.Elem()
		}

		if f.Kind() == reflect.Struct {
			a := structFields(f)
			if len(a) > 0 {
				m = append(m, a...)
			}

			continue
		}

		if f.Kind() == reflect.Slice && f.Len() > 0 {
			for x := 0; x < f.Len(); x++ {
				a := structFields(f.Index(x))
				if len(a) > 0 {
					m = append(m, a...)
				}
			}

			continue
		}

		s, ok := v.Type().Field(i).Tag.Lookup("required")
		if !ok {
			continue
		}

		if isEmpty(f) || f.Kind() == reflect.Invalid {
			msg := Message{
				Index:  i,
				Field:  getFieldName(v.Type().Field(i)),
				ErrMsg: "this field is required",
			}

			if s != "-" {
				msg.ErrMsg = s
			}

			m = append(m, msg)
		}
	}

	return m
}

// isEmpty check if the field value is empty
func isEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0

	case reflect.Float32, reflect.Float64:
		return v.Float() == 0

	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0

	case reflect.String:
		return v.String() == ""

	case reflect.Map:
		return len(v.MapKeys()) == 0

	case reflect.Slice:
		return v.Len() == 0
	}

	return false
}

// getFieldName get the field name to use in some API response
// this way, the struct can be easy encoded to xml or json
func getFieldName(f reflect.StructField) (s string) {
	if t, ok := f.Tag.Lookup("json"); ok {
		s = t
	} else if t, ok := f.Tag.Lookup("xml"); ok {
		s = t
	} else {
		s = f.Name
	}

	return
}
