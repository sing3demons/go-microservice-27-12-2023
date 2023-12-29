package store

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sing3demons/category/utils"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Store struct {
	Client *mongo.Client
	DB     *gorm.DB
}

func NewStore() *Store {
	mongoUri := os.Getenv("MONGO_URL")
	mysqlDsn := os.Getenv("MYSQL_DSN")
	postgresDsn := os.Getenv("POSTGRES_DSN")
	dsn := os.Getenv("DSN")

	if mongoUri != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoUri))
		if err != nil {
			panic(err)
		}

		if err := client.Ping(ctx, readpref.Primary()); err != nil {
			panic(err)
		}

		log.WithFields(log.Fields{
			"URI":  utils.EncryptAES(mongoUri),
			"TYPE": "MONGO",
		}).Info("Connected to MongoDB!")

		return &Store{Client: client}
	} else if mysqlDsn != "" {
		db, err := gorm.Open(mysql.Open(mysqlDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Connected to MySQL!")
		return &Store{DB: db}
	} else if postgresDsn != "" {
		db, err := gorm.Open(postgres.Open(postgresDsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}

		fmt.Println("Connected to MySQL!")
		return &Store{DB: db}
	} else {
		if dsn == "" {
			dsn = "gorm.db"
		}
		db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		return &Store{DB: db}
	}
}
