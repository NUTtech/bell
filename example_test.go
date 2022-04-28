package bell_test

import (
	"fmt"
	"github.com/nuttech/bell"
	"sort"
)

type CustomStruct struct {
	name  string
	param int32
}

func Example() {
	event := "event_name"
	event2 := "event_name_2"

	// add listener on event event_name
	bell.Listen(event, func(message bell.Message) {
		// we extend CustomStruct in message
		customStruct := message.(CustomStruct)
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

	// wait until the event completes its work
	bell.Wait()

	// Output:
	// [event_name event_name_2]
	// [event_name]
	// false
	// {testName 12}
}
