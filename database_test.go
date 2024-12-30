package pocketsql

import "testing"

func TestConnectionInfo(t *testing.T) {
	driverName, dataSourceName := connectionInfo("sqlite:file:database")
	if driverName != "sqlite" {
		t.Errorf(`Expected 'sqlite', got '%s'`, driverName)
	}
	if dataSourceName != "file:database" {
		t.Errorf(`Expected 'file:database', got '%s'`, dataSourceName)
	}

	driverName, dataSourceName = connectionInfo("postgres:user=pqgotest dbname=pqgotest sslmode=verify-full")
	if driverName != "postgres" {
		t.Errorf(`Expected 'postgres', got '%s'`, driverName)
	}
	if dataSourceName != "user=pqgotest dbname=pqgotest sslmode=verify-full" {
		t.Errorf(`Expected 'user=pqgotest dbname=pqgotest sslmode=verify-full', got '%s'`, dataSourceName)
	}

	driverName, dataSourceName = connectionInfo("postgres://user:password@localhost/db?sslmode=verify-full")
	if driverName != "postgres" {
		t.Errorf(`Expected 'postgres', got '%s'`, driverName)
	}
	if dataSourceName != "postgres://user:password@localhost/db?sslmode=verify-full" {
		t.Errorf(`Expected 'postgres://user:password@localhost/db?sslmode=verify-full', got '%s'`, dataSourceName)
	}
}
