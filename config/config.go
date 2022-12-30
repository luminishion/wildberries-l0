package config

const (
	DatabaseURL = "postgres://postgres:gogo@host.docker.internal:5432/postgres"

	WebAddr = ":8080"

	NatsClusterID = "test-cluster"
	NatsClientID  = "test-service"
	NatsChannel   = "orders"
	NatsURL       = "nats://host.docker.internal:4222"
)
