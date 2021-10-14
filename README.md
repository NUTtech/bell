# Bell

Bell is the simplest event system written in Go (Golang) which is based on the execution of handlers independent of the main channel.

- Support for custom event data.
- Internally, it launches each handler in a separate goroutine and passes messages to them through channels, the handlers are executed independently of the main thread.
- Support for adding multiple handlers to an event.
- Complete unit testing.

## Installation

To install Bell package, you need install [Go](https://golang.org) with [modules](https://github.com/golang/go/wiki/Modules) support and set Go workspace first.
1. Use the below Go command to install Bell:
```shell
go get -u github.com/nuttech/bell
```
2. Import package in your code:
```go
import "github.com/nuttech/bell"
```

## Examples
```go
import (
	"fmt"
	"github.com/nuttech/bell"
	"sort"
	"time"
)

type CustomStruct struct {
	name string
	param int32
}

func Example() {
	event := "event_name"
	event2 := "event_name_2"

	// add listener on event event_name
	bell.Listen(event, func(message bell.Message) {
		// we extend CustomStruct in message.Value
		customStruct := message.Value.(CustomStruct)
		fmt.Println(customStruct)
	})
	// add listener on event event_name_2
	bell.Listen(event2, func(message bell.Message) {

	})

	// get event list
	list := bell.List()

	// only for test
	sort.Strings(list)
	fmt.Println(list)

	// remove listeners on event_name_2
	bell.Remove(event2)

	// get event list again
	fmt.Println(bell.List())

	// check if exists event_name_2 event in storage
	fmt.Println(bell.Has(event2))

	// call event event_name
	_ = bell.Ring(event, CustomStruct{name: "testName", param: 12})

	// ONLY FOR EXAMPLE
	// add sleep because the event handler does not have time
	// to be processed before the completion of the script execution
	time.Sleep(time.Millisecond * 50)

	// Output:
	// [event_name event_name_2]
	// [event_name]
	// false
	// {testName 12}
}
```
