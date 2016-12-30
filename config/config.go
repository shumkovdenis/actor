package config

type accountAPI struct {
	URL  string
	Type string
}

type ratesAPI struct {
	URL     string
	Timeout int
}

type updateServer struct {
	URL            string
	UpdateInterval int
}

type config struct {
	Version      string
	AccountAPI   accountAPI
	RatesAPI     ratesAPI
	UpdateServer updateServer
}

// Conf server
var Conf = config{
	Version: "0.0.1",
	AccountAPI: accountAPI{
		URL:  "https://silver-pay.com/auth/clientjson",
		Type: "BINOPT",
	},
	RatesAPI: ratesAPI{
		URL:     "http://currency.silver-pay.com:3030/rates",
		Timeout: 5 * 1000,
	},
	UpdateServer: updateServer{
		URL:            "http://localhost:8080",
		UpdateInterval: 5 * 1000,
	},
}
