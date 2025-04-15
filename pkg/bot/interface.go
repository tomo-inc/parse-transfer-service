package bot

import "context"

type Level int32

const (
	Info Level = iota
	Warning
	Error
)

type KeyValue struct {
	Key   string
	Value interface{}
}

type Msg struct {
	Title string
	Level Level
	Data  []*KeyValue
}

type Notificator interface {
	Send(context.Context, Msg) error
}
