package initializers

import (
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/rgraterol/movies-reviewers-wrapper/config"
)

var DatabaseConfig DatabaseConfiguration

// DatabaseConfiguration represents a database configuration.
type DatabaseConfiguration struct {
	// URL is the database address.
	URL string `yaml:"url"`
	// MaxIdleConns sets the maximum number of connections in the idle connection pool.
	MaxIdleConns int `yaml:"maxIdleConns"`
	// MaxOpenConns sets the maximum number of open connections to the database.
	MaxOpenConns int `yaml:"maxOpenConns"`
	// ConnMaxLifetime sets the maximum amount of time in minutes a connection may be reused.
	ConnMaxLifetime int `yaml:"connMaxLifetime"`
	// Automigrate set condition to automatically migrate db schema.
	AutoMigrate bool `yaml:"autoMigrate"`
}

func DatabaseInitializer() {
	err := LoadConfigSection("database", &DatabaseConfig)
	if err != nil {
		panic(errors.Wrap(err, "failed to read the database config"))
	}

	if url := os.Getenv("DATABASE_URL"); url != "" {
		DatabaseConfig.URL = url
	}

	config.DB, err = gorm.Open(mysql.Open(DatabaseConfig.URL), &gorm.Config{Logger: initGormLogger()})
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize the DB"))
	}
	pool, err := config.DB.DB()
	if err != nil {
		panic(errors.Wrap(err, "failed to configure connection pool"))
	}
	pool.SetMaxIdleConns(DatabaseConfig.MaxIdleConns)
	pool.SetMaxOpenConns(DatabaseConfig.MaxOpenConns)
	pool.SetConnMaxLifetime(time.Duration(DatabaseConfig.ConnMaxLifetime))
}

func MockDatabaseInitializer() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: initGormLogger()})
	if err != nil {
		panic(errors.Wrap(err, "failed to connect gorm with mock DB"))
	}
}


func initGormLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:              time.Second,
			LogLevel:                   logger.Silent,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}
