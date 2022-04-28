# Bell

[![GoDoc](https://pkg.go.dev/badge/github.com/nuttech/bell?status.svg)](https://pkg.go.dev/github.com/nuttech/bell?tab=doc)
[![Release](https://img.shields.io/github/release/nuttech/bell.svg?style=flat)](https://github.com/nuttech/bell/releases)
[![codecov](https://codecov.io/gh/NUTtech/bell/branch/master/graph/badge.svg?token=3TMnbQkEny)](https://codecov.io/gh/NUTtech/bell)
[![Tests](https://github.com/NUTtech/bell/actions/workflows/tests.yml/badge.svg)](https://github.com/NUTtech/bell/actions/workflows/tests.yml)

Bell is the simplest event system written in Go (Golang) which is based on the execution of handlers independent of the
main channel.

- Written in pure go. Has no third-party libraries.
- Support for custom event data.
- Internally, it launches each handler in a separate goroutine and passes messages to them through channels, the
  handlers are executed independently of the main thread.
- Support for adding multiple handlers to an event.
- Complete unit testing.

## Installation

To install Bell package, you need install [Go](https://golang.org)
with [modules](https://github.com/golang/go/wiki/Modules) support and set Go workspace first.

1. Use the below Go command to install Bell:

```shell
go get -u github.com/nuttech/bell
```

2. Import package in your code:

```go
import "github.com/nuttech/bell"
```

## Usage

### Adding event listener

The handler function accepts the Message type as input

```go
bell.Listen("event_name", func(message bell.Message) {
	// here you must write your handler code
})
```

`bell.Message` has `interface{}` type and can consist any data.

You can add more handlers one event:

```go
bell.Listen("event_name", func(message bell.Message) { 
	// first handler
})
bell.Listen("event_name", func(message bell.Message) {
	// second handler
})
```

### Calling an event

This code call event. Activating handlers, who subscribed on "event_name" event

```go
bell.Call("event_name", "some data")

bell.Call("event_name", 1) // int

bell.Call("event_name", false) // bool
```

If you passing struct type of data:

```go
type userStruct struct {
	Name string
}
bell.Call("event_name", userStruct{Name: "Jon"})
```

Then parsing the data in the handler may look like this:

```go
bell.Listen("event_name", func(message bell.Message) {
	user := message.(userStruct)
	
	fmt.Printf("%#v\n", userStruct{Name: "Jon"})  // main.userStruct{Name:"Jon"}
})
```

### Getting events list

To get a list of events to which handlers are subscribed, call the code:

```go
bell.List()
```

### Checking if exists listeners of event

You can check the existence of subscribers to an event like this:

```go
bell.Has("event_name")
```

### Removing listeners of event (all events)

You can delete all listeners or listeners of only one event.

#### Removing all listeners on all events

```go
_ = bell.Remove()
```

#### Removing listeners of only the event "event_name"

```go
_ = bell.Remove("event_name")
```

### Wait until all events finish their work

```go
bell.Wait()
```

### Change events queue size (apply only for new listeners)

```go
bell.Queue(42)
```

### Usage without global state

You can also use the bell package without using global state. To do this, you need to create a state storage object
and use it.

```go
events := bell.New()
events.Listen("event", func(message bell.Message) {})
_ = events.Ring("event", "Hello bell!")
```

## Examples

See full example in [example_test.go](example_test.go).
