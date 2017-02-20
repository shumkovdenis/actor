package core

type Command interface {
	Command() string
}

type Event interface {
	Event() string
}

type Code interface {
	Code() string
}
