package config

type KafkaConfig struct {
	Brokers []string
	Topic   string
	GroupID string
}

func NewKafkaConfig() KafkaConfig {
	return KafkaConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
		GroupID: "my-group",
	}
}
