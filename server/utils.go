package server

type Tell func(interface{})

type Proc func(Tell) bool

func Process(proc Proc, tell Tell) bool {
	var r bool

	ch := make(chan interface{})

	t := func(m interface{}) {
		ch <- m
	}

	go func() {
		r = proc(t)
		close(ch)
	}()

	for m := range ch {
		tell(m)
	}

	return r
}
