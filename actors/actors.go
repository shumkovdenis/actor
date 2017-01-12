package actors

type Tell func(interface{})

type Proc func(Tell)

func Process(proc Proc, tell Tell) {
	ch := make(chan interface{})

	t := func(m interface{}) {
		ch <- m
	}

	go func() {
		proc(t)
		close(ch)
	}()

	for m := range ch {
		tell(m)
	}
}
