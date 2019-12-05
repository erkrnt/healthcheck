package main

import (
	"flag"
	"os"

	"github.com/docker/docker/client"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

// Flags : command-line flags
type Flags struct {
	Debug   bool
	Preseed bool
}

// Service : service configuration
type Service struct {
	DB     *gorm.DB
	Docker *client.Client
	Env    string
	Flags
}

// loadDB : Loads database from configuration
func loadDB(env string, flags Flags) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", "data.db")
	if err != nil {
		return nil, err
	}
	if flags.Debug == true {
		db.LogMode(true)
	}
	db.AutoMigrate(&HealthcheckType{})
	if flags.Preseed == true {
		preseed(db)
	}
	return db, nil
}

// loadDocker : Loads new Docker client to use.
func loadDocker() (*client.Client, error) {
	cli, err := client.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return cli, nil
}

// loadEnv : See docs https://github.com/joho/godotenv#precendence--conventions
func loadEnv() string {
	env := os.Getenv("HEALTHCHECK_ENV")
	if "" == env {
		env = "development"
	}
	godotenv.Load(".env." + env + ".local")
	if "test" != env {
		godotenv.Load(".env.local")
	}
	godotenv.Load(".env." + env)
	godotenv.Load()
	return env
}

// loadFlags : loads command-line flags
func loadFlags() Flags {
	debug := flag.Bool("debug", false, "Enables debug output.")
	preseed := flag.Bool("preseed", false, "Enables preseeding of database for required types.")
	flag.Parse()
	return Flags{Debug: *debug, Preseed: *preseed}
}

// preseed : loads tables with preseed values (types, etc)
func preseed(db *gorm.DB) {
	db.FirstOrCreate(&HealthcheckType{}, HealthcheckType{Name: "ping"})
}

// Initialize : loads app statup functions
func Initialize() (*Service, error) {
	flags := loadFlags()
	env := loadEnv()
	db, err := loadDB(env, flags)
	if err != nil {
		return nil, err
	}
	docker, err := loadDocker()
	if err != nil {
		return nil, err
	}
	return &Service{DB: db, Docker: docker, Env: env, Flags: flags}, nil
}
