package responder

type respOpts struct {
	statusCode int
}

type ResponseOption func(*respOpts)

func WithStatusCode(code int) ResponseOption {
	return func(o *respOpts) {
		o.statusCode = code
	}
}
