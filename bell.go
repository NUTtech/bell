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

// Message The message that is passed to the event handler
type Message interface{}

// Events thread safe structure stores events, their handlers and functions for management
type Events struct {
	sync.RWMutex
	channels  map[string][]chan Message
	wg        sync.WaitGroup
	queueSize uint
}

// New constructor for Events
func New() *Events {
	return &Events{channels: map[string][]chan Message{}}
}

// Queue set events queue size
func (e *Events) Queue(size uint) *Events {
	e.Lock()
	defer e.Unlock()

	e.queueSize = size
	return e
}

// ListenN Subscribe on event where
// event - the event name,
// handlerFunc - handler function
// copiesCount - count handlers copies run
func (e *Events) ListenN(event string, handlerFunc func(message Message), copiesCount uint) {
	e.Lock()
	defer e.Unlock()

	channel := make(chan Message, e.queueSize)

	for i := uint(0); i < copiesCount; i++ {
		go func(c chan Message, wg *sync.WaitGroup) {
			for {
				message, ok := <-c
				if !ok {
					break
				}
				handlerFunc(message)
				wg.Done()
			}
		}(channel, &e.wg)
	}

	e.channels[event] = append(e.channels[event], channel)
}

// Listen Subscribe on event where
// event - the event name,
// handlerFunc - handler function
func (e *Events) Listen(event string, handlerFunc func(message Message)) {
	e.ListenN(event, handlerFunc, 1)
}

// Ring Call event there
// event - event name
// message - data that will be passed to the event handler
func (e *Events) Ring(event string, message Message) error {
	e.RLock()
	defer e.RUnlock()

	if _, ok := e.channels[event]; !ok {
		return fmt.Errorf("channel %s not found", event)
	}

	for _, c := range e.channels[event] {
		e.wg.Add(1)
		c <- message
	}
	return nil
}

// Has Checks if there are listeners for the passed event
func (e *Events) Has(event string) bool {
	e.RLock()
	defer e.RUnlock()

	_, ok := e.channels[event]
	return ok
}

// List Returns a list of events that listeners are subscribed to
func (e *Events) List() []string {
	e.RLock()
	defer e.RUnlock()

	list := make([]string, 0, len(e.channels))
	for event := range e.channels {
		list = append(list, event)
	}
	return list
}

// Remove Removes listeners by event name
// Removing listeners closes channels and stops the goroutine.
//
// If you call the function without the "names" parameter, all listeners of all events will be removed.
func (e *Events) Remove(names ...string) {
	e.Lock()
	defer e.Unlock()

	if len(names) == 0 {
		keys := make([]string, 0, len(e.channels))
		for k := range e.channels {
			keys = append(keys, k)
		}

		names = keys
	}

	for _, name := range names {
		for _, channel := range e.channels[name] {
			close(channel)
		}

		delete(e.channels, name)
	}
}

// Wait Blocks the thread until all running events are completed
func (e *Events) Wait() {
	e.Lock()
	defer e.Unlock()

	e.wg.Wait()
}

// globalState store of global event handlers
var globalState = New()

// ListenN Subscribe on event where
// event - the event name,
// handlerFunc - handler function
// copiesCount - count handlers copies run
func ListenN(event string, handlerFunc func(message Message), copiesCount uint) {
	globalState.ListenN(event, handlerFunc, copiesCount)
}

// Listen Subscribe on event where
// event - the event name,
// handlerFunc - handler function
func Listen(event string, handlerFunc func(message Message)) {
	globalState.Listen(event, handlerFunc)
}

// Ring Call event there
// event - event name
// message - data that will be passed to the event handler
func Ring(event string, message Message) error {
	return globalState.Ring(event, message)
}

// Has Checks if there are listeners for the passed event
func Has(event string) bool {
	return globalState.Has(event)
}

// List Returns a list of events that listeners are subscribed to
func List() []string {
	return globalState.List()
}

// Remove Removes listeners by event name
// Removing listeners closes channels and stops the goroutine.
//
// If you call the function without the "names" parameter, all listeners of all events will be removed.
func Remove(names ...string) {
	globalState.Remove(names...)
}

// Wait Blocks the thread until all running events are completed
func Wait() {
	globalState.Wait()
}

// Queue set events queue size
func Queue(size uint) {
	globalState.Queue(size)
}
