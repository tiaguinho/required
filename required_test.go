package required

import (
	"fmt"
	"testing"
)

//
type S struct {
	StringDefault string `required:"-"`
	StringCustom  string `required:"wow! this is empty"`
	IntDefault    int32  `required:"-"`
	IntCustom     int32  `required:"where is the number?"`
}

var (
	v = S{}

	stringCusto = "wow! this is empty"
)

//
func TestValidate(t *testing.T) {
	if msgs, err := Validate(v); err != nil {
		fmt.Println(msgs)
	}
}
