package impl_example

// Logger определяет интерфейс для логирования
type Logger interface {
	Log(msg string)
	Logf(format string, args ...interface{})
	Error(msg string)
	Errorf(format string, args ...interface{})
}
