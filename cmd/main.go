package main

import (
	"abir"
	"abir/pkg/config"
	"abir/pkg/handler"
	"abir/pkg/repository"
	"abir/pkg/repository/postgres"
	"abir/pkg/service"
	"abir/pkg/storage"
	"flag"
	"github.com/streadway/amqp"
	//ps "abir/proto/proto"
	"context"
	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go"
	_ "github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

var rootPath *string

func init() {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	rootPath = flag.String("root_path", "", "Root path")
	flag.Parse()

	logFilePath, err := filepath.Abs(*rootPath + "logs/api.log")
	if err != nil {
		logrus.Fatal(err)
	}
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)
	logrus.SetOutput(mw)

	if err := config.Init(*rootPath); err != nil {
		logrus.Fatalf("error loading config: %s\n", err.Error())
	}
	envFilePath, err := filepath.Abs(*rootPath + ".env")
	if err != nil {
		logrus.Fatal(err)
	}
	if err := godotenv.Load(envFilePath); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}
}

func main() {
	//res, err := utils.GetCardToken("8600140209489880", "2511")
	//if err != nil {
	//	logrus.Fatalf("%s", err)
	//}
	//logrus.Print(res.Card.Phone)
	conn, err := amqp.Dial(os.Getenv("RABBIT_MQ"))
	if err != nil {
		logrus.Fatalf("Error occurred on ampq connection: %s\n", err.Error())
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatalf("Error occurred on ampq channel: %s\n", err.Error())
	}
	defer ch.Close()

	minioStorage, err := initStorage()
	if err != nil {
		logrus.Fatalf("Error occurred on storage initialization: %s\n", err.Error())
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	_, err = redisClient.Ping().Result()
	if err != nil {
		logrus.Fatalf("failed to connect redis: %s", err.Error())
	}

	dashboardDb, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
		Schema:   viper.GetString("db.dashboard_schema"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize dashboard db: %s", err.Error())
	}
	publicDb, err := postgres.NewPostgresDB(postgres.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
		Schema:   viper.GetString("db.public_schema"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize public db: %s", err.Error())
	}
	repos := repository.NewRepository(dashboardDb, publicDb)
	services := service.NewService(repos, redisClient, minioStorage, ch)
	handlers := handler.NewHandler(services)

	srv := abir.NewServer()
	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running server: %s", err.Error())
		}
	}()

	logrus.Print("Abir started")

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	logrus.Print("Abir shutting down")
	_, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	if err = dashboardDb.Close(); err != nil {
		logrus.Errorf("error occurred while closing db connection: %s\n", err.Error())
	}
	if err = publicDb.Close(); err != nil {
		logrus.Errorf("error occurred while closing db connection: %s\n", err.Error())
	}
}

func initStorage() (storage.Storage, error) {
	client, err := minio.New(
		viper.GetString("storage.url"),
		os.Getenv("ACCESS_KEY"),
		os.Getenv("SECRET_KEY"), false)
	if err != nil {
		return nil, err
	}

	exists, err := client.BucketExists(viper.GetString("storage.bucket"))
	if err != nil {
		return nil, err
	}

	logrus.Infof("Bucket %s exists: %v", viper.GetString("storage.bucket"), exists)

	return storage.NewFileStorage(client,
		viper.GetString("storage.bucket"),
		viper.GetString("storage.url"),
		os.Getenv("HOST")), nil
}
