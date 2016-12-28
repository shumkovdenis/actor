package config

type accountAPI struct {
	URL  string
	Type string
}

type ratesAPI struct {
	URL     string
	Timeout int
}

type config struct {
	AccountAPI accountAPI
	RatesAPI   ratesAPI
}

// Conf server
var Conf = config{
	AccountAPI: accountAPI{
		URL:  "http://stage.silver-pay.com/auth/clientjson",
		Type: "BINOPT",
	},
	RatesAPI: ratesAPI{
		URL:     "http://currency.silver-pay.com:3030/rates",
		Timeout: 5000,
	},
}
