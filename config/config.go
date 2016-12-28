package config

type accountAPI struct {
	URL  string
	Type string
}

type config struct {
	AccountAPI accountAPI
}

// Conf server
var Conf = config{
	AccountAPI: accountAPI{
		URL:  "http://stage.silver-pay.com/auth/clientjson",
		Type: "BINOPT",
	},
}
