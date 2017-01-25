package session

// Login -> command.login
type Login struct {
}

// LoginSuccess -> event.login.success
type LoginSuccess struct {
	// Client string `json:"client"`
}

// LoginFail -> event.login.fail
type LoginFail struct {
	Message string `json:"message"`
}

/*
// Join -> command.join
type Join struct {
	Client string `mapstructure:"client"`
}

// JoinSuccess -> event.join.success
type JoinSuccess struct {
}

// JoinFail -> event.join.fail
type JoinFail struct {
	Message string `json:"message"`
}
*/

// Fail -> event.session.fail
type Fail struct {
	Message string `json:"message"`
}
