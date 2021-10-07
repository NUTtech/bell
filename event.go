// Package event реализует простую систему событиый
//
// На каждое событие (event) можно добавить несколько обработчиков (handlerFunc)
// Вызов обработчиков события происходит в отдельной горутине через установленный канал
// При вызове определенного события происходит последовательная передача сообщения (Message)
// каждому обработчику события.
//
// Если канал закрывается, горутина для этого события прекращает свою работу.
//
// Пример использования:
// On("event_name", func(message Message) { fmt.PrintLn(message) }) - добавляем обработчик события event_name
// Call("event_name", "some_data") - Вызоваем события "event_name", тем самым запуская обработчики установленные ранее
package event

import (
	"fmt"
	"sync"
	"time"
)

// Хранилище состояний обработчиков событий
var eventMap = &events{channels: map[string][]chan Message{}}

// Message Сообщение которое передается в обработчик события
type Message struct {
	Event     string
	Timestamp time.Time
	Value     interface{}
}

type events struct {
	sync.RWMutex
	channels map[string][]chan Message
}

// On Добавление обработчика события
// event - название/код события
// handlerFunc - функция-обработчик события. На вход принимает структуру Message
func On(event string, handlerFunc func(message Message)) {
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

// Call Вызывает событие
// event - название/код события
// value - данные, которые будут переданы в функции-обработчики события внутри Message
func Call(event string, value interface{}) error {
	eventMap.RLock()
	defer eventMap.RUnlock()

	if _, ok := eventMap.channels[event]; !ok {
		return fmt.Errorf("channel %s not found", event)
	}

	for _, c := range eventMap.channels[event] {
		c <- Message{Event: event, Timestamp: time.Now(), Value: value}
	}
	return nil
}

// Has Возвращает true если существуют обработчики переданного события
func Has(event string) bool {
	eventMap.RLock()
	defer eventMap.RUnlock()

	_, ok := eventMap.channels[event]
	return ok
}

// List Возвращает список событий, на которые установлены обработчики
func List() []string {
	eventMap.RLock()
	defer eventMap.RUnlock()

	var list []string
	for event := range eventMap.channels {
		list = append(list, event)
	}
	return list
}

// Remove Удаляет обработчики события или событий
// При удалении обработчиков закрываются каналы и прекращают работу горутины
//
// Если вызвать функцию без параметра names - будут удалены все обработчики всех событий
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
