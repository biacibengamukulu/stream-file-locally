package migrate

import (
	"fmt"
	"strconv"
	"strings"

	"enterprise-delivery-platform/shared/cassandra"
)

// EnsureKeyspace creates the keyspace if it doesn't exist yet.
// rf is the replication factor (SimpleStrategy) typically "1" for dev/single node.
func EnsureKeyspace(system *cassandra.Session, keyspace string, rf string) error {
	if system == nil || system.Session == nil {
		return fmt.Errorf("system session is nil")
	}
	keyspace = strings.TrimSpace(keyspace)
	if keyspace == "" {
		return fmt.Errorf("keyspace is empty")
	}

	rfi := 1
	if strings.TrimSpace(rf) != "" {
		if n, err := strconv.Atoi(strings.TrimSpace(rf)); err == nil && n > 0 {
			rfi = n
		}
	}

	cql := fmt.Sprintf(
		"CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': %d}",
		keyspace, rfi,
	)
	return system.Session.Query(cql).Exec()
}

// Run executes the given CQL statements against the provided session.
// It is safe to run multiple times if your DDL uses IF NOT EXISTS.
func Run(sess *cassandra.Session, ddls []string) error {
	if sess == nil || sess.Session == nil {
		return fmt.Errorf("session is nil")
	}
	for i, stmt := range ddls {
		s := strings.TrimSpace(stmt)
		if s == "" {
			continue
		}
		// gocql doesn't need trailing semicolons
		s = strings.TrimSuffix(s, ";")
		if err := sess.Session.Query(s).Exec(); err != nil {
			return fmt.Errorf("migration %d failed: %w (cql=%q)", i, err, s)
		}
	}
	return nil
}
