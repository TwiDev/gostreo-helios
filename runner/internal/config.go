package internal

import "os"

type Config struct {
	GRPCAddr     string
	WSAddr       string
	TemporalHost string
	RabbitURL    string
	PostgresDSN  string
}

func LoadConfig() Config {
	return Config{
		GRPCAddr:     getEnv("GRPC_ADDR", ":50051"),
		WSAddr:       getEnv("WS_ADDR", ":8080"),
		TemporalHost: getEnv("TEMPORAL_HOST", "localhost:7233"),
		RabbitURL:    getEnv("RABBIT_URL", "amqp://guest:guest@rabbitmq:5672/"),
		PostgresDSN:  getEnv("POSTGRES_DSN", "postgres://postgres:password@postgres:5432/workflows?sslmode=disable"),
	}
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
