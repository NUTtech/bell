package event

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

// resetSystem Очистка хранилища состояний обработчиков событий
func resetSystem() {
	for k := range eventMap.channels {
		for _, channel := range eventMap.channels[k] {
			close(channel)
		}
	}
	eventMap = &events{channels: map[string][]chan Message{}}
}

// TestOn Проверка работы функции добавления обработчика событий
func TestOn(t *testing.T) {
	resetSystem()
	defer resetSystem()

	expMessageEvent := "test_event"
	expMessageValue := "value"

	On(expMessageEvent, func(message Message) {
		assert.Equal(t, expMessageEvent, message.Event)
		assert.Equal(t, expMessageValue, message.Value)
	})

	assert.Equal(t, 1, len(eventMap.channels))
	assert.Equal(t, 1, len(eventMap.channels[expMessageEvent]))

	assert.NotPanics(t, func() {
		err := Call(expMessageEvent, expMessageValue)
		assert.NoError(t, err)
	})

	assert.NotPanics(t, func() {
		resetSystem()
	})
}

// TestCall_Fail Проверка корректности ошибки при ошибочном вызове события
func TestCall_Fail(t *testing.T) {
	resetSystem()
	defer resetSystem()

	err := Call("undefined_event", func() {})
	assert.EqualError(t, err, "channel undefined_event not found")
}

// TestRemove проверка удаления обработчиков события из хранилища
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

// TestHas Проверка корректности определения существования обработчиков события
func TestHas(t *testing.T) {
	resetSystem()
	defer resetSystem()

	assert.False(t, Has("test"))

	eventMap.channels["test"] = append(eventMap.channels["test"], make(chan Message))
	assert.True(t, Has("test"))
}

// TestList Проверка корректного получения списка событий, на которые установлены обработчики
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
