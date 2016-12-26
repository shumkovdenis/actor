package account

type Balance struct {
}

type BalanceSuccess struct {
}

func balance(msg *Balance) interface{} {
	return &BalanceSuccess{}
}
