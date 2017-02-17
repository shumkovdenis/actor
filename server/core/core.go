package core

type Fail interface {
	Code() string
}

func IsFail(msg interface{}) bool {
	_, ok := msg.(Fail)
	return ok
}
