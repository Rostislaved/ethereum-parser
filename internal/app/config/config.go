package config

type Config struct {
	Provider Provider
	Parser   Parser
	Server   Server
}

func Get() Config {
	configProvider := Provider{
		URL:                 "https://cloudflare-eth.com",
		ClientTimeoutInSecs: 5,
	}

	configParser := Parser{
		IntervalInSecs:            5,
		InitialBlockNumber:        16050110,
		NumberOfFetchingWorkers:   5,
		NumberOfProcessingWorkers: 100,
		NumberOfSavingWorkers:     10,
	}

	configServer := Server{
		Addr:            ":8080",
		ShutdownTimeout: 2,
	}

	return Config{
		Provider: configProvider,
		Parser:   configParser,
		Server:   configServer,
	}
}

type Provider struct {
	URL                 string
	ClientTimeoutInSecs int
}

type Parser struct {
	IntervalInSecs            int
	InitialBlockNumber        int64
	NumberOfFetchingWorkers   int
	NumberOfProcessingWorkers int
	NumberOfSavingWorkers     int
}

type Server struct {
	Addr            string
	ShutdownTimeout int
}
