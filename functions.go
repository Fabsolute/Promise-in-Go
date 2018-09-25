package promise

func New(resolver func(resolve, reject func(interface{}))) *Promise {
	p := &Promise{pending, nil, make([]*Handler, 0)}
	go doResolve(resolver, p.resolve, p.reject)
	return p
}

func FromFunction(fun func() interface{}) *Promise {
	return New(func(resolve, reject func(interface{})) {
		if fun == nil {
			reject("fun cannot be nil")
			return
		}

		defer func() {
			err := recover()
			if err != nil {
				reject(err)
			}
		}()
		response := fun()
		resolve(response)
	})
}

func Resolve(value interface{}) *Promise {
	return New(func(resolve func(interface{}), reject func(interface{})) {
		resolve(value)
	})
}

func Reject(reason interface{}) *Promise {
	return New(func(resolve func(interface{}), reject func(interface{})) {
		reject(reason)
	})
}

func All(promises ...*Promise) *Promise {
	accumulator := make([]interface{}, len(promises))
	ready := Resolve(nil)

	for i, promise := range promises {
		i := i
		promise := promise
		ready = ready.Then(func(_ interface{}) interface{} {
			return promise
		}, nil).Then(func(value interface{}) interface{} {
			accumulator[i] = value
			return nil
		}, nil)
	}

	return ready.Then(func(_ interface{}) interface{} {
		return accumulator
	}, nil)
}

func Race(promises ...*Promise) *Promise {
	return New(func(resolve func(interface{}), reject func(interface{})) {
		for _, promise := range promises {
			promise := promise
			promise.Then(func(value interface{}) interface{} {
				resolve(value)
				return value
			}, func(reason interface{}) interface{} {
				reject(reason)
				return reason
			})
		}
	})
}

func doResolve(fn func(_, _ func(value interface{})), onFulfilled, onRejected func(value interface{})) {
	done := false
	defer func() {
		err := recover()
		if err != nil {
			if done {
				return
			}
			done = true
			onRejected(err)
		}
	}()

	fn(func(value interface{}) {
		if done {
			return
		}
		done = true
		onFulfilled(value)
	}, func(reason interface{}) {
		if done {
			return
		}
		done = true
		onRejected(reason)
	})
}

func getThen(value interface{}) (func(onFulfilled, onRejected func(reason interface{})), bool) {
	promise, ok := value.(Promise)
	if ok {
		return func(onFulfilled, onRejected func(reason interface{})) {
			resolve := func(value interface{}) interface{} {
				if onFulfilled != nil {
					onFulfilled(value)
				}
				return nil
			}

			reject := func(value interface{}) interface{} {
				if onRejected != nil {
					onRejected(value)
				}
				return nil
			}
			promise.Then(resolve, reject)
		}, true
	}

	return nil, false
}

func isPromise(value interface{}) bool {
	_, ok := value.(*Promise)
	return ok
}
