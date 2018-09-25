package promise

type Handler struct {
	onFulfilled func(value interface{})
	onRejected  func(reason interface{})
}

func NewHandler(onFulfilled, onRejected func(value interface{})) *Handler {
	return &Handler{onFulfilled: onFulfilled, onRejected: onRejected}
}
