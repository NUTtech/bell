// Package bell implements a simple event system (bell ringing and listening)
//
// Several listeners can be added for each ringing (handlerFunc).
// Listeners are called in a separate goroutine through an established channel.
// When the bell rings, a message is sequentially transmitted to each listener.
//
// If a channel is closed, the goroutine for that event is terminated.
//
// Example for usage:
// Listen("event_name", func(message Message) { fmt.PrintLn(message) }) - add listener on bell by name "event_name"
// Ring("event_name", "some_data") - Ring on bell (call event "event_name")
package bell

import (
	"fmt"
	"sync"
)

// State store of event handlers
var eventMap = &events{channels: map[string][]chan Message{}}

// Message The message that is passed to the event handler
type Message interface{}

type events struct {
	sync.RWMutex
	channels map[string][]chan Message
}

// Listen Subscribe on event where
// event - the event name,
// handlerFunc - handler function
func Listen(event string, handlerFunc func(message Message)) {
	eventMap.Lock()
	defer eventMap.Unlock()

	channel := make(chan Message)

	go func(c chan Message) {
		for {
			message, ok := <-c
			if !ok {
				break
			}

			handlerFunc(message)
		}
	}(channel)

	eventMap.channels[event] = append(eventMap.channels[event], channel)
}

// Ring Call event there
// event - event name
// message - data that will be passed to the event handler
func Ring(event string, message Message) error {
	eventMap.RLock()
	defer eventMap.RUnlock()

	if _, ok := eventMap.channels[event]; !ok {
		return fmt.Errorf("channel %s not found", event)
	}

	for _, c := range eventMap.channels[event] {
		c <- message
	}
	return nil
}

// Has Checks if there are listeners for the passed event
func Has(event string) bool {
	eventMap.RLock()
	defer eventMap.RUnlock()

	_, ok := eventMap.channels[event]
	return ok
}

// List Returns a list of events that listeners are subscribed to
func List() []string {
	eventMap.RLock()
	defer eventMap.RUnlock()

	var list []string
	for event := range eventMap.channels {
		list = append(list, event)
	}
	return list
}

// Remove Removes listeners by event name
// Removing listeners closes channels and stops the goroutine.
//
// If you call the function without the "names" parameter, all listeners of all events will be removed.
func Remove(names ...string) {
	eventMap.Lock()
	defer eventMap.Unlock()

	if len(names) == 0 {
		keys := make([]string, 0, len(eventMap.channels))
		for k := range eventMap.channels {
			keys = append(keys, k)
		}

		names = keys
	}

	for _, name := range names {
		for _, channel := range eventMap.channels[name] {
			close(channel)
		}

		delete(eventMap.channels, name)
	}
}
