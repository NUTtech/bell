package bell

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// resetSystem Clearing the State Store of Event Listeners
func resetSystem() {
	for k := range eventMap.channels {
		for _, channel := range eventMap.channels[k] {
			close(channel)
		}
	}
	eventMap = &events{channels: map[string][]chan Message{}}
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

	assert.Equal(t, 1, len(eventMap.channels))
	assert.Equal(t, 1, len(eventMap.channels[expMessageEvent]))

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

	eventMap.channels["test"] = append(eventMap.channels["test"], make(chan Message), make(chan Message))
	eventMap.channels["test2"] = append(eventMap.channels["test2"], make(chan Message))

	Remove("test")
	assert.Equal(t, 1, len(eventMap.channels))

	eventMap.channels["test3"] = append(eventMap.channels["test3"], make(chan Message))
	eventMap.channels["test4"] = append(eventMap.channels["test4"], make(chan Message))
	Remove("test2")
	assert.Equal(t, 2, len(eventMap.channels))

	eventMap.channels["test3"] = append(eventMap.channels["test3"], make(chan Message))
	eventMap.channels["test4"] = append(eventMap.channels["test4"], make(chan Message))
	Remove()
	assert.Equal(t, 0, len(eventMap.channels))
}

// TestHas Checking the Correctness of Determining the Existence of Event Listeners
func TestHas(t *testing.T) {
	resetSystem()
	defer resetSystem()

	assert.False(t, Has("test"))

	eventMap.channels["test"] = append(eventMap.channels["test"], make(chan Message))
	assert.True(t, Has("test"))
}

// TestList Checking the correct receipt of the list of events on which handlers are installed
func TestList(t *testing.T) {
	resetSystem()
	defer resetSystem()

	assert.Empty(t, List())

	eventMap.channels["test"] = append(eventMap.channels["test"], make(chan Message), make(chan Message))
	eventMap.channels["test2"] = append(eventMap.channels["test2"], make(chan Message))

	actualList := List()
	sort.Strings(actualList)

	assert.Equal(t, 2, len(actualList))
	assert.Equal(t, []string{"test", "test2"}, actualList)
}
