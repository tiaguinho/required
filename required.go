package required

import (
	"fmt"
	"reflect"
)

//
type Message struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

//
func Validate(v interface{}) (messages []Message, err error) {
	numFields := reflect.ValueOf(v).NumField()
	if numFields > 0 {
		messages = make([]Message, 0)

		for i := 0; i < numFields; i++ {
			if msg, ok := reflect.TypeOf(v).Field(i).Tag.Lookup("required"); ok {
				fmt.Println(msg)
			}
		}
	}

	return
}
