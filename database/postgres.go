package database

import (
	"fmt"
	"log"
	"sync"

	"github.com/ganyacc/Ganesh_OrderProcessingSystem/config"
	"github.com/ganyacc/Ganesh_OrderProcessingSystem/entities"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type postgresDatabase struct {
	Db *gorm.DB
}

var (
	once       sync.Once
	DbInstance *postgresDatabase
)

// NewPostgresDatabase returns the new instance of postgres db
func NewPostgresDatabase(conf *config.Config) Database {
	once.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			conf.Db.Host,
			conf.Db.User,
			conf.Db.Password,
			conf.Db.DBName,
			conf.Db.Port,
			conf.Db.SSLMode,
			conf.Db.TimeZone,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic("failed to connect database")
		}

		logrus.Printf("connected to '%v' database", conf.Db.DBName)

		DbInstance = &postgresDatabase{Db: db}
	})

	return DbInstance
}

// CloseDb closes the underlying db connection
func (p *postgresDatabase) CloseDb(db *gorm.DB) error {
	sqlDb, err := db.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}

func (p *postgresDatabase) AutoMigrateTables() error {

	//--Create the recurring enum type if it does not exist
	err := p.Db.Exec(`
		CREATE TYPE IF NOT EXISTS order_status AS ENUM ('unfulfilled', 'fulfilled');
	`).Error
	if err != nil {
		log.Println("Error creating order_status enum:", err)
		return err
	}

	return p.Db.AutoMigrate(&entities.Customer{}, &entities.Product{}, &entities.Order{})

}

func (p *postgresDatabase) GetDb() *gorm.DB {
	return DbInstance.Db
}
