package data

import (
	"flag"
	"fmt"
	"strings"
	"sync"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	log "github.com/sirupsen/logrus"

	"github.com/jinzhu/gorm"
)

var (
	db     *gorm.DB
	onceDB sync.Once
)

type DBParams struct {
	Name         string
	Host         string
	Port         int
	User         string
	Password     string
	Retries      int
	Log          bool
	MaxOpenConns int
	MaxIdleConns int
}

func (dbParam *DBParams) DeclFlags() {
	flag.StringVar(&dbParam.Name, "db.name", "ultra", "database name")
	flag.StringVar(&dbParam.Host, "db.host", "127.0.0.1", "database host")
	flag.IntVar(&dbParam.Port, "db.port", 5432, "database port")
	flag.StringVar(&dbParam.User, "db.user", "postgres", "database username")
	flag.StringVar(&dbParam.Password, "db.password", "", "database password")
	flag.IntVar(&dbParam.Retries, "db.retries", 3,
		"number of retries to access database")
	flag.BoolVar(&dbParam.Log, "db.log", true, "true if log db transactions")
	flag.IntVar(&dbParam.MaxOpenConns, "db.max_open_conn", 16,
		"max number of open connections")
	flag.IntVar(&dbParam.MaxIdleConns, "db.max_idle_conn", 16,
		"max number of idle connections")
}

func Open(params *DBParams) {
	onceDB.Do(func() {
		dbParam := map[string]interface{}{
			"dbname":   params.Name,
			"host":     params.Host,
			"port":     params.Port,
			"user":     params.User,
			"password": params.Password,
			"sslmode":  "disable",
		}
		dbStrings := make([]string, 0, len(dbParam))
		for dbKey, dbVal := range dbParam {
			if dbStr, ok := dbVal.(string); !ok || len(dbStr) > 0 {
				dbStrings = append(dbStrings,
					fmt.Sprintf("%s=%v", dbKey, dbVal))
			}
		}
		dbstring := strings.Join(dbStrings, " ")
		log.Info(dbstring)

		var err error
		for i := 0; i <= params.Retries; i++ {
			log.Infof("Attempt to connect to database: %s:%d", params.Host, params.Port)
			db, err = gorm.Open("postgres", dbstring)
			if err == nil {
				log.Infof("Successfully connected to a database.")
				break
			} else {
				// Exponential back-off.
				waitSeconds := 1 << uint(i)
				log.Warnf("Exponential back-off: %ds", waitSeconds)
				time.Sleep(time.Duration(waitSeconds) * time.Second)
			}
		}

		if err != nil {
			log.Errorf("Number of attempts exceeded maximum, giving up.")
			log.Fatalf(err.Error())
		}

		db.DB().SetMaxOpenConns(params.MaxOpenConns)
		db.DB().SetMaxIdleConns(params.MaxIdleConns)
		db.LogMode(params.Log)
	})
}

func Close() {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func DB() *gorm.DB {
	return db
}

func Transaction(f func(db *gorm.DB) error) (err error) {
	transaction := DB().Begin()
	if err = transaction.Error; err != nil {
		return
	}

	if err = f(transaction); err != nil {
		if err = transaction.Rollback().Error; err != nil {
			log.Errorf("Failed to rollback DB transaction: %v", err)
		}
		return err
	}

	return transaction.Commit().Error
}
