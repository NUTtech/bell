package bell

import (
	"sort"
	"sync/atomic"
	"testing"
	"time"
)

// resetSystem Clearing the State Store of Event Listeners
func resetSystem() {
	for k := range globalState.channels {
		for _, channel := range globalState.channels[k] {
			close(channel)
		}
	}
	globalState = &Events{channels: map[string][]chan Message{}}
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Error(err)
	}
}

func assertTrue(t *testing.T, v bool) {
	t.Helper()
	if v != true {
		t.Error("Value must be true")
	}
}

// TestListenN checking the function of adding multiple copies of event listeners
func TestListenN(t *testing.T) {
	defer resetSystem()

	eventName := "event"
	var wasRunning int32
	ListenN(eventName, func(Message) { atomic.AddInt32(&wasRunning, 1) }, 3)

	assertNoError(t, Ring(eventName, nil))
	Wait()

	assertTrue(t, wasRunning == 1)

	assertNoError(t, Ring(eventName, nil))
	assertNoError(t, Ring(eventName, nil))
	Wait()

	assertTrue(t, wasRunning == 3)
}

// TestListen Testing the function of adding event listeners
func TestListen(t *testing.T) {
	defer resetSystem()

	expMessageEvent := "test_event"
	expMessageValue := "value"

	Listen(expMessageEvent, func(message Message) {
		assertTrue(t, expMessageValue == message)
	})

	assertTrue(t, len(globalState.channels) == 1)
	assertTrue(t, len(globalState.channels[expMessageEvent]) == 1)

	assertNoError(t, Ring(expMessageEvent, expMessageValue))
}

// TestRing_Fail Checking the correctness of error handling in case of an erroneous ringing
func TestRing_Fail(t *testing.T) {
	defer resetSystem()

	err := Ring("undefined_event", func() {})
	assertTrue(t, err.Error() == "channel undefined_event not found")
}

// TestRemove Checking if event handlers are removed from storage
func TestRemove(t *testing.T) {
	defer resetSystem()

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message), make(chan Message))
	globalState.channels["test2"] = append(globalState.channels["test2"], make(chan Message))

	Remove("test")
	assertTrue(t, len(globalState.channels) == 1)

	globalState.channels["test3"] = append(globalState.channels["test3"], make(chan Message))
	globalState.channels["test4"] = append(globalState.channels["test4"], make(chan Message))
	Remove("test2")
	assertTrue(t, len(globalState.channels) == 2)

	globalState.channels["test3"] = append(globalState.channels["test3"], make(chan Message))
	globalState.channels["test4"] = append(globalState.channels["test4"], make(chan Message))
	Remove()
	assertTrue(t, len(globalState.channels) == 0)
}

// TestHas Checking the Correctness of Determining the Existence of Event Listeners
func TestHas(t *testing.T) {
	defer resetSystem()

	assertTrue(t, !Has("test"))

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message))
	assertTrue(t, Has("test"))
}

// TestList Checking the correct receipt of the list of events on which handlers are installed
func TestList(t *testing.T) {
	defer resetSystem()

	assertTrue(t, len(List()) == 0)

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message), make(chan Message))
	globalState.channels["test2"] = append(globalState.channels["test2"], make(chan Message))

	actualList := List()
	sort.Strings(actualList)

	assertTrue(t, len(actualList) == 2)
	assertTrue(t, actualList[0] == "test")
	assertTrue(t, actualList[1] == "test2")
}

// TestWait Checking Wait function
func TestWait(t *testing.T) {
	defer resetSystem()

	eventName := "test"
	var wasRunning int32

	Listen(eventName, func(Message) {
		time.Sleep(time.Millisecond)
		atomic.StoreInt32(&wasRunning, 1)
	})
	assertNoError(t, Ring(eventName, nil))

	Wait()

	assertTrue(t, wasRunning == 1)
}

// TestQueue checking function for set queue size
func TestQueue(t *testing.T) {
	defer resetSystem()

	var size uint = 6
	Queue(size)
	assertTrue(t, size == globalState.queueSize)
}
