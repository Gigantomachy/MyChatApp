package db

import "github.com/gocql/gocql"

func InitSession(hosts []string, keyspace, username, password string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.RoundRobinHostPolicy())
	if username != "" { // local dev/tests run auth-less Cassandra
		cluster.Authenticator = gocql.PasswordAuthenticator{Username: username, Password: password}
	}
	return cluster.CreateSession()
}
