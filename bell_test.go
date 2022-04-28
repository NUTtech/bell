package bell

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

// TestListen Testing the function of adding event listeners
func TestListen(t *testing.T) {
	resetSystem()
	defer resetSystem()

	expMessageEvent := "test_event"
	expMessageValue := "value"

	Listen(expMessageEvent, func(message Message) {
		assert.Equal(t, expMessageValue, message)
	})

	assert.Equal(t, 1, len(globalState.channels))
	assert.Equal(t, 1, len(globalState.channels[expMessageEvent]))

	assert.NotPanics(t, func() {
		err := Ring(expMessageEvent, expMessageValue)
		assert.NoError(t, err)
	})

	assert.NotPanics(t, func() {
		resetSystem()
	})
}

// TestRing_Fail Checking the correctness of error handling in case of an erroneous ringing
func TestRing_Fail(t *testing.T) {
	resetSystem()
	defer resetSystem()

	err := Ring("undefined_event", func() {})
	assert.EqualError(t, err, "channel undefined_event not found")
}

// TestRemove Checking if event handlers are removed from storage
func TestRemove(t *testing.T) {
	resetSystem()
	defer resetSystem()

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message), make(chan Message))
	globalState.channels["test2"] = append(globalState.channels["test2"], make(chan Message))

	Remove("test")
	assert.Equal(t, 1, len(globalState.channels))

	globalState.channels["test3"] = append(globalState.channels["test3"], make(chan Message))
	globalState.channels["test4"] = append(globalState.channels["test4"], make(chan Message))
	Remove("test2")
	assert.Equal(t, 2, len(globalState.channels))

	globalState.channels["test3"] = append(globalState.channels["test3"], make(chan Message))
	globalState.channels["test4"] = append(globalState.channels["test4"], make(chan Message))
	Remove()
	assert.Equal(t, 0, len(globalState.channels))
}

// TestHas Checking the Correctness of Determining the Existence of Event Listeners
func TestHas(t *testing.T) {
	resetSystem()
	defer resetSystem()

	assert.False(t, Has("test"))

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message))
	assert.True(t, Has("test"))
}

// TestList Checking the correct receipt of the list of events on which handlers are installed
func TestList(t *testing.T) {
	resetSystem()
	defer resetSystem()

	assert.Empty(t, List())

	globalState.channels["test"] = append(globalState.channels["test"], make(chan Message), make(chan Message))
	globalState.channels["test2"] = append(globalState.channels["test2"], make(chan Message))

	actualList := List()
	sort.Strings(actualList)

	assert.Equal(t, 2, len(actualList))
	assert.Equal(t, []string{"test", "test2"}, actualList)
}

// TestWait Checking Wait function
func TestWait(t *testing.T) {
	resetSystem()
	defer resetSystem()

	eventName := "test"
	var wasRunning int32

	Listen(eventName, func(Message) {
		time.Sleep(time.Millisecond)
		atomic.StoreInt32(&wasRunning, 1)
	})
	require.NoError(t, Ring(eventName, nil))

	Wait()

	assert.Equal(t, int32(1), wasRunning)
}

// TestQueue checking function for set queue size
func TestQueue(t *testing.T) {
	var size uint = 6
	Queue(size)
	assert.Equal(t, size, globalState.queueSize)
}
