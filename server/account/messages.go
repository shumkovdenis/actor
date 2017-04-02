package account

type Message interface {
	AccountMessage()
}

type Authorize struct {
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (*Authorize) AccountMessage() {}

func (*Authorize) Command() string {
	return "command.account.authorize"
}

type Authorized struct {
	Categories []Category `json:"categories"`
}

func (*Authorized) AccountMessage() {}

func (*Authorized) Event() string {
	return "event.account.authorized"
}

type AlreadyAuthorized struct{}

func (*AlreadyAuthorized) AccountMessage() {}

func (*AlreadyAuthorized) Event() string {
	return "event.account.already_authorized"
}

type NotAuthorized struct{}

func (*NotAuthorized) AccountMessage() {}

func (*NotAuthorized) Event() string {
	return "event.account.not_authorized"
}

type AuthorizationFailed struct{}

func (*AuthorizationFailed) AccountMessage() {}

func (m *AuthorizationFailed) Event() string {
	return "event.account.authorize.failed"
}

type GetBalance struct{}

func (*GetBalance) AccountMessage() {}

func (*GetBalance) Command() string {
	return "command.account.balance"
}

type Balance struct {
	Balance float64 `json:"balance"`
}

func (*Balance) AccountMessage() {}

func (*Balance) Event() string {
	return "event.account.balance"
}

type GetBalanceFailed struct{}

func (*GetBalanceFailed) AccountMessage() {}

func (*GetBalanceFailed) Event() string {
	return "event.account.balance.failed"
}

type GetGameSession struct {
	GameID string `mapstructure:"game_id"`
}

func (*GetGameSession) AccountMessage() {}

func (*GetGameSession) Command() string {
	return "command.account.session"
}

type GameSession struct {
	SessionID string `json:"session_id"`
	GameID    string `json:"game_id"`
	ServerURL string `json:"server_url"`
}

func (*GameSession) AccountMessage() {}

func (*GameSession) Event() string {
	return "event.account.session"
}

type GetGameSessionFailed struct{}

func (*GetGameSessionFailed) AccountMessage() {}

func (*GetGameSessionFailed) Event() string {
	return "event.account.session.failed"
}

type Withdraw struct{}

func (*Withdraw) AccountMessage() {}

func (*Withdraw) Command() string {
	return "command.account.withdraw"
}

type WithdrawSuccess struct{}

func (*WithdrawSuccess) AccountMessage() {}

func (*WithdrawSuccess) Event() string {
	return "event.account.withdraw"
}

type WithdrawFailed struct{}

func (*WithdrawFailed) AccountMessage() {}

func (*WithdrawFailed) Event() string {
	return "event.account.withdraw.failed"
}

type Cashback struct{}

func (*Cashback) AccountMessage() {}

func (*Cashback) Command() string {
	return "command.account.cashback"
}

type CashbackSuccess struct {
	Cashback float64 `json:"cashback"`
}

func (*CashbackSuccess) AccountMessage() {}

func (*CashbackSuccess) Event() string {
	return "event.account.cashback"
}

type CashbackFailed struct{}

func (*CashbackFailed) AccountMessage() {}

func (*CashbackFailed) Event() string {
	return "event.account.cashback.failed"
}

type GetJackpotsTops struct{}

func (*GetJackpotsTops) AccountMessage() {}

func (*GetJackpotsTops) Command() string {
	return "command.account.jackpots.tops"
}

type JackpotsTops struct {
	Tops []Jackpot `json:"tops"`
}

func (*JackpotsTops) AccountMessage() {}

func (*JackpotsTops) Event() string {
	return "event.account.jackpots.tops"
}

type GetJackpotsTopsFailed struct{}

func (*GetJackpotsTopsFailed) AccountMessage() {}

func (*GetJackpotsTopsFailed) Event() string {
	return "event.account.jackpots.tops.failed"
}

type GetJackpotsList struct{}

func (*GetJackpotsList) AccountMessage() {}

func (*GetJackpotsList) Command() string {
	return "command.account.jackpots.list"
}

type JackpotsList struct {
	Large  string `json:"large"`
	Medium string `json:"medium"`
	Small  string `json:"small"`
}

func (*JackpotsList) AccountMessage() {}

func (*JackpotsList) Event() string {
	return "event.account.jackpots.list"
}

type GetJackpotsListFailed struct{}

func (*GetJackpotsListFailed) AccountMessage() {}

func (*GetJackpotsListFailed) Event() string {
	return "event.account.jackpots.list.failed"
}

type StartLiveJackpotsTops struct{}

func (*StartLiveJackpotsTops) AccountMessage() {}

type StopLiveJackpotsTops struct{}

func (*StopLiveJackpotsTops) AccountMessage() {}

type StartLiveJackpotsList struct{}

func (*StartLiveJackpotsList) AccountMessage() {}

type StopLiveJackpotsList struct{}

func (*StopLiveJackpotsList) AccountMessage() {}

type Category struct {
	Title string `json:"title"`
	Games []Game `json:"games"`
}

type Game struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Thumb string `json:"thumb"`
}

type Jackpot struct {
	Account string `json:"account"`
	Win     string `json:"win"`
	Date    string `json:"date"`
}
