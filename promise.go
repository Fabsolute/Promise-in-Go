package promise

import (
	"sync"
)

const (
	pending = iota
	fulfilled
	rejected
)

type Promise struct {
	state    int
	value    interface{}
	handlers []*Handler
}

func (promise *Promise) Then(onFulfilled func(value interface{}) interface{}) *Promise {
	return promise.thenInternal(onFulfilled, nil)
}

func (promise *Promise) Catch(onRejected func(reason interface{}) interface{}) *Promise {
	return promise.thenInternal(nil, onRejected)
}

func (promise *Promise) Await() interface{} {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	var result interface{} = nil
	success := false
	promise.thenInternal(func(value interface{}) interface{} {
		result = value
		success = true
		return value
	}, func(reason interface{}) interface{} {
		result = reason
		success = false
		return reason
	}).Then(func(value interface{}) interface{} {
		wg.Done()
		return value
	})
	wg.Wait()

	if !success {
		panic(result)
	}

	return result
}


func (promise *Promise) done(onFulfilled func(value interface{}), onRejected func(reason interface{})) {
	handler := NewHandler(onFulfilled, onRejected)
	go promise.executeHandler(handler)
}

func (promise *Promise) thenInternal(onFulfilled func(value interface{}) interface{}, onRejected func(reason interface{}) interface{}) *Promise {
	return New(func(resolve, reject func(interface{})) {
		defer func() {
			err := recover()
			if err != nil {
				reject(err)
			}
		}()
		promise.done(func(result interface{}) {
			makeFulfillChain(result, onFulfilled, onRejected, resolve, reject)
		}, func(reason interface{}) {
			makeRejectChain(reason, onFulfilled, onRejected, resolve, reject)
		})
	})
}

func (promise *Promise) fulfill(value interface{}) {
	promise.state = fulfilled
	promise.value = value
	promise.executeHandlers()
}

func (promise *Promise) reject(reason interface{}) {
	promise.state = rejected
	promise.value = reason
	promise.executeHandlers()
}

func (promise *Promise) resolve(value interface{}) {
	defer promise.handlePanic()
	then, ok := getThen(value)
	if ok {
		doResolve(then, promise.resolve, promise.reject)
		return
	}

	promise.fulfill(value)
}

func (promise *Promise) executeHandlers() {
	for _, handler := range promise.handlers {
		handler := handler
		promise.executeHandler(handler)
	}

	promise.handlers = make([]*Handler, 0)
}

func (promise *Promise) executeHandler(handler *Handler) {
	if promise.state == pending {
		promise.handlers = append(promise.handlers, handler)
	} else {
		if promise.state == fulfilled && handler.onFulfilled != nil {
			handler.onFulfilled(promise.value)
		}
		if promise.state == rejected && handler.onRejected != nil {
			handler.onRejected(promise.value)
		}
	}
}

func (promise *Promise) handlePanic() {
	err := recover()
	if err != nil {
		promise.reject(err)
	}
}
