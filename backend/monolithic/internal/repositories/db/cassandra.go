package db

import "github.com/gocql/gocql"

func InitSession(hosts []string, keyspace string) (*gocql.Session, error) {
	cluster := gocql.NewCluster(hosts...)
	cluster.Keyspace = keyspace
	cluster.Consistency = gocql.Quorum
	return cluster.CreateSession()
}