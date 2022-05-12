package mock

import (
	"fmt"

	"github.com/ngocphuongnb/tetua/app/logger"
)

type MockLoggerMessage struct {
	Type   string
	Params []interface{}
}

type MockLogger struct {
	Silence  bool
	Messages []*MockLoggerMessage
}

func (l *MockLogger) WithContext(context logger.Context) logger.Logger {
	return l
}

func (l *MockLogger) Last() MockLoggerMessage {
	if len(l.Messages) == 0 {
		return MockLoggerMessage{}
	}
	return *l.Messages[len(l.Messages)-1]
}

func (l *MockLogger) Info(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Info", Params: params})
}

func (l *MockLogger) Debug(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Debug", Params: params})
}
func (l *MockLogger) Warn(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Warn", Params: params})
}
func (l *MockLogger) Error(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Error", Params: params})
}
func (l *MockLogger) DPanic(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "DPanic", Params: params})
}
func (l *MockLogger) Panic(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Panic", Params: params})
}
func (l *MockLogger) Fatal(params ...interface{}) {
	if !l.Silence {
		fmt.Println(params...)
	}
	l.Messages = append(l.Messages, &MockLoggerMessage{Type: "Fatal", Params: params})
}
