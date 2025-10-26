package http

type HTTPServer interface {
	Start() error
	Stop() error
}
