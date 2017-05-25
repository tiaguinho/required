package required

import (
	"reflect"
	"fmt"
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
	Field   string `json:"field" xml:"field"`
	ErrMsg  string `json:"message" xml:"message"`
}

// Validate return error if any field is left empty
func Validate(v interface{}) error {
	sm := checkFields(reflect.ValueOf(v))
	if len(sm) == 0 {
		return nil
	}
	return &errorMessages{sm}
}

// ValidateWithMessage return error and slice of message if any field is left empty
func ValidateWithMessage(v interface{}) (error, []Message) {
	sm := checkFields(reflect.ValueOf(v))
	if len(sm) == 0 {
		return nil, sm
	}
	return &errorMessages{sm}, sm
}

// checkFields check the type of field and if the field has a required tag
func checkFields(v reflect.Value) []Message {
	sm := make([]Message, 0)
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Ptr && v.Field(i).Kind() != reflect.Struct {
			if s, ok := v.Type().Field(i).Tag.Lookup("required"); ok {
				if isEmpty(v.Field(i)) {
					m := Message{
						Field:   getFieldName(v.Type().Field(i)),
						ErrMsg: "this field is required",
					}

					if s != "-" {
						m.ErrMsg = s
					}

					sm = append(sm, m)
				}
			}
		} else {
			a := checkFields(v.Field(i))
			if len(a) > 0 {
				sm = append(sm, a...)
			}
		}
	}

	return sm
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