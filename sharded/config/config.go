package config

import (
	"os"
	"strings"

	"github.com/biacibengamukulu/stream-file-locally/sharded/constants"
)

type Config struct {
	ServiceName string

	HTTPPort string

	CassandraHosts    string
	CassandraKeyspace string

	KafkaBrokers string
	KafkaGroupID string

	DropboxRefreshToken string
}

func Getenv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func sanitizeServiceName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, " ", "_")
	return s
}

func Load(serviceName string) Config {
	svc := sanitizeServiceName(serviceName)

	defaultKeyspace := constants.CASSANDRA_KEYSPACE_PREFIX + svc
	// Examples:
	//  - delivery   => biatechlibs_stream

	return Config{
		ServiceName: serviceName,
		HTTPPort:    Getenv("HTTP_PORT", "8080"),

		CassandraHosts:    Getenv("CASSANDRA_HOSTS", constants.CASSANDRA_HOSTS),
		CassandraKeyspace: Getenv("CASSANDRA_KEYSPACE", defaultKeyspace),

		KafkaBrokers:        Getenv("KAFKA_BROKERS", constants.KAFKA_HOST),
		KafkaGroupID:        Getenv("KAFKA_GROUP_ID", svc+"-group"),
		DropboxRefreshToken: Getenv("DROPBOX_REFRESH_TOKEN", constants.DROPBOX_REFRESH_TOKEN),
	}
}
