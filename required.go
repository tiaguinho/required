package required

import (
	"reflect"
	"fmt"
)

//
type Message struct {
	Field   string `json:"field" xml:"field"`
	Message string `json:"message" xml:"message"`
}

//
func Validate(v interface{}) ([]Message, error) {
	return checkFields(reflect.ValueOf(v))
}

//
func checkFields(v reflect.Value) (messages []Message, err error) {
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Kind() != reflect.Ptr && v.Field(i).Kind() != reflect.Struct {
			if msg, ok := v.Type().Field(i).Tag.Lookup("required"); ok {
				if isEmpty(v.Field(i)) {
					message := Message{
						Field:   getFieldName(v.Type().Field(i)),
						Message: "this field is required",
					}

					if msg != "-" {
						message.Message = msg
					}

					messages = append(messages, message)
				}
			}
		} else {
			s, err := checkFields(v.Field(i))
			if err != nil {
				messages = append(messages, s...)
			}
		}
	}

	if len(messages) > 0 {
		err = fmt.Errorf("%d fields are required", len(messages))
	}
	
	return
}

//
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

//
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