package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/pprof"

	"github.com/chenmuyao/secumon/internal/event/monitor"
	"github.com/chenmuyao/secumon/internal/repository"
	"github.com/chenmuyao/secumon/internal/repository/cache"
	"github.com/chenmuyao/secumon/internal/repository/dao"
	"github.com/chenmuyao/secumon/internal/service"
	svc "github.com/chenmuyao/secumon/internal/service/logmonitor"
	"github.com/chenmuyao/secumon/internal/web/logmonitor"
	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Test struct {
	gorm.Model
}

func main() {
	// slog.SetLogLoggerLevel(slog.LevelDebug)

	// Init DB (GORM)
	db := InitDB()
	// Init Redis
	redis := InitRedis()
	// Init RabbitMQ
	amqpConn := InitRabbitMQ()
	defer amqpConn.Close()

	// Init Services

	// Run Webserver
	server := gin.Default()
	server.GET("/", func(ctx *gin.Context) {
		// Test redis
		err := redis.Set(ctx, "test", "secumon", time.Hour).Err()
		if err != nil {
			panic(err)
		}

		// Test rabbit
		ch, err := amqpConn.Channel()
		if err != nil {
			panic(err)
		}

		q, err := ch.QueueDeclare("secumon", false, false, false, false, nil)
		if err != nil {
			panic(err)
		}

		// Sender
		body := "test"
		err = ch.PublishWithContext(ctx, "", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
		if err != nil {
			panic(err)
		}

		// Receiver
		msgs, err := ch.ConsumeWithContext(ctx, q.Name, "", true, false, false, false, nil)
		if err != nil {
			panic(err)
		}
		for msg := range msgs {
			log.Println("read:", string(msg.Body))
			break
		}
		ctx.String(http.StatusOK, "So far so good")
	})

	const exchangeName = "api-security-logs"
	qName, err := monitor.AccessLogMQSetup(amqpConn, exchangeName)
	if err != nil {
		panic(err)
	}
	publisher := monitor.NewRabbitMQLogMonitorPublisher(amqpConn, exchangeName)

	logDAO := dao.NewLogDAO(db)

	bfChecker := cache.NewBruteForceChecker(redis)
	htChecker := cache.NewHighTrafficChecker(redis)
	alertCache := cache.NewRedisAlertCache(redis)

	logRepo := repository.NewLogRepo(logDAO)
	alertRepo := repository.NewAlertRepo(logDAO, alertCache)

	bfDetector := svc.NewBruteForceDetector(logRepo, bfChecker, alertCache)
	htDetector := svc.NewHighTrafficDetector(logRepo, htChecker, alertCache)
	alertSvc := service.NewAlertService(alertRepo)

	err = monitor.NewRabbitMQLogMonitorConsumer(amqpConn).
		UseBruteForceDetector(bfDetector).
		UseHighTrafficDetector(htDetector).
		StartConsumer(exchangeName, qName)
	if err != nil {
		panic(err)
	}

	hdl := logmonitor.NewLogHandler(publisher, alertSvc)
	hdl.RegisterHandlers(server)

	pprof.Register(server)
	server.Run(":8989")
}

func InitDB() *gorm.DB {
	dsn := "host=localhost user=postgres password=postgres dbname=secumon port=15432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	// Test DB
	err = db.AutoMigrate(&Test{}, &dao.SecurityEvent{})
	if err != nil {
		panic(err)
	}
	return db
}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:16379",
	})
}

func InitRabbitMQ() *amqp.Connection {
	connection, err := amqp.Dial("amqp://secumon:secumon@localhost:5672")
	if err != nil {
		panic(err)
	}
	return connection
}
