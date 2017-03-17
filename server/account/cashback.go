package account

func cashback(username, password string) Message {
	return &CashbackSuccess{500.00}
}
