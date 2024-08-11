package scoringconf

type Config struct {
	ServerPort      string             `json:"server_port"`
	Concurrent      Concurrent         `json:"concurrent,omitempty"`
	Redis           Redis              `json:"redis"`
	Persistence     *PersistentStorage `json:"persistence,omitempty"`
	EventServiceUrl string             `json:"event_service_url"`
}

type Concurrent struct {
	NumberOfWorker int `json:"number_of_worker,omitempty"`
	MaxQueueLength int `json:"max_queue_length,omitempty"`
}

type Redis struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type PersistentStorage struct {
	Enabled bool   `json:"enabled,omitempty"`
	URI     string `json:"uri,omitempty"`
}
