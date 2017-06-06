# Required [![Build Status](https://travis-ci.org/tiaguinho/required.svg?branch=master)](https://travis-ci.org/tiaguinho/required) [![GoDoc](https://godoc.org/github.com/tiaguinho/required?status.png)](https://godoc.org/github.com/tiaguinho/required) [![Go Report Card](https://goreportcard.com/badge/github.com/tiaguinho/required)](https://goreportcard.com/report/github.com/tiaguinho/required)
Required is used to validate if the field value is empty

### Install

```bash
go get github.com/tiaguinho/required
```

### Example

Simplest way to use the package:

```go
package main

import (
    "github.com/tiaguinho/required"
    "log"
)

type Test struct {
    FirstName string `json:"first_name" required:"-"`
    LastName  string `json:"last_name" required:"last name is required"`
}

func main()  {
    t := Test{
        FirstName: "Go",
        LastName: "Required",
    }
    
    if err := required.Validate(t); err != nil {
        log.Println(err)
    }
}
```

If you like to get all the validation messages with field's name to return in some API, just change to this:

```go
func main() {
    t := Test{
         FirstName: "Go",
         LastName: "Required",
     }
     
     if msg, err := required.ValidateWithMessage(t); err != nil {
         log.Println(err, msg)
     }
}
```
