package log

type NopLogger struct{}

func (NopLogger) Printf(string, ...any) {}

func NewNopLogger() Logger {
    return NopLogger{}
}