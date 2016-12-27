package app

type Update struct {
}

type Restart struct {
}

func update(msg *Update) interface{} {
	return &Restart{}
}
