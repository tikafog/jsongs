# JSONGS

## This project is forked from `encoding.json`

### Added method fetching for private variables because I couldn't stand the crappy rules of the original package

## Usage:

- All usage is the same as the official package `encoding.json`

## Add Features

```go
package main

import (
	"fmt"

	jsongs "github.com/tikafog/jsongs"
)

type Example struct {
	name string `json:"name" json-getter:"MyName" json-setter:"SetMyName"`
}

func (receiver Example) MyName() string {
	return receiver.name
}

func (receiver *Example) SetMyName(name string) {
	receiver.name = name
}

func main() {
	v, err := jsongs.Marshal(&Example{
		name: "my name is jsongs",
	})
	if err != nil {
		panic(err)
	}
	//dosomething
	fmt.Println(string(v))
	// Output:
	//{"name":"my name is jsongs"}
}
```

- The getter/setter tags for internal variables are optional,
  getter defaults to camel case:`Name`, setter defaults to Set + camel case:`SetName`

