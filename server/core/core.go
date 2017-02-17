package core

type Command interface {
	Command() string
}

type Event interface {
	Event() string
}

type Fail interface {
	Fail() string
}
