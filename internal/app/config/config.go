package config

type Parser struct {
	IntervalInSecs            int
	InitialBlockNumber        int64
	NumberOfFetchingWorkers   int
	NumberOfProcessingWorkers int
	NumberOfSavingWorkers     int
}
