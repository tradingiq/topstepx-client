package services

type noopLogger struct{}

func (n noopLogger) Log(keyvals ...interface{}) error {
	return nil
}

func newNoopLogger() noopLogger {
	return noopLogger{}
}