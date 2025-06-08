package postgres

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dwilanang/psp/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Postgres driver
)

var (
	masterInstances = map[string]*sqlx.DB{}
	slaveInstances  = map[string]*sqlx.DB{}
	mu              sync.RWMutex
)

type Config struct {
	Driver       string
	MasterDSN    string
	SlaveDSN     string
	MaxOpenConns int
	MaxIdleConns int
	Timeout      int
}

func initDB(name string, cfg Config) error {
	mu.Lock()
	defer mu.Unlock()
	if _, exists := masterInstances[name]; exists {
		return nil
	}

	// Init master
	masterDB, err := sqlx.Connect(cfg.Driver, cfg.MasterDSN)
	if err != nil {
		return err
	}
	masterDB.SetConnMaxLifetime(time.Hour)
	masterDB.SetMaxOpenConns(cfg.MaxOpenConns)
	masterDB.SetMaxIdleConns(cfg.MaxIdleConns)
	if cfg.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
		defer cancel()
		if err := masterDB.PingContext(ctx); err != nil {
			return err
		}
	}
	masterInstances[name] = masterDB

	// Init slave (optional)
	if cfg.SlaveDSN != "" {
		slaveDB, err := sqlx.Connect(cfg.Driver, cfg.SlaveDSN)
		if err != nil {
			return err
		}
		slaveDB.SetConnMaxLifetime(time.Hour)
		slaveDB.SetMaxOpenConns(cfg.MaxOpenConns)
		slaveDB.SetMaxIdleConns(cfg.MaxIdleConns)
		if cfg.Timeout > 0 {
			ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Timeout)*time.Second)
			defer cancel()
			if err := slaveDB.PingContext(ctx); err != nil {
				return err
			}
		}
		slaveInstances[name] = slaveDB
	}

	return nil
}

func getMasterDB(name string) *sqlx.DB {
	mu.RLock()
	defer mu.RUnlock()
	db, ok := masterInstances[name]
	if !ok {
		log.Fatalf("Master DB instance %s not found, call InitDB first", name)
	}
	return db
}

func getSlaveDB(name string) *sqlx.DB {
	mu.RLock()
	defer mu.RUnlock()
	db, ok := slaveInstances[name]
	if !ok {
		log.Fatalf("Slave DB instance %s not found, call InitDB first with SlaveDSN", name)
	}
	return db
}

func healthCheckDB(name string) error {
	db := getMasterDB(name)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return db.PingContext(ctx)
}

func setupDB(cfg *config.Config) {

	dbConfig := Config{
		Driver:       cfg.DBDriver,
		MasterDSN:    fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", cfg.DBUser, cfg.DBPassword, cfg.DBName),
		SlaveDSN:     "",
		MaxOpenConns: 10,
		MaxIdleConns: 5,
		Timeout:      30,
	}
	err := initDB("main", dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

func Connect(cfg *config.Config) *sqlx.DB {
	// Initialize the database connection
	setupDB(cfg)

	// Check the health of the Postgres database
	checkPostgresHealth()

	// Get the master database connection
	dbConn := getMasterDB("main")

	return dbConn
}

func checkPostgresHealth() {
	err := healthCheckDB("main")
	if err != nil {
		log.Printf("Postgres health check failed: %v", err)
	} else {
		log.Println("Postgres is healthy")
	}
}
