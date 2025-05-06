package utils

import (
	"context"
	"fmt"
	"ginchat/config"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

// 使用viper加载配置
func LoadConfig() (*config.Database, error) {
	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("config")
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config: %v", err)
	}
	var db config.Database
	if err := viper.UnmarshalKey("database", &db); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %v", err)
	}
	return &db, nil
}

// 初始化mysql数据库连接
func InitDB(cfg *config.MySQLConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
	)
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // 禁用彩色打印
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	//连接池配置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	if timeout, err := time.ParseDuration(cfg.Timeout); err != nil {
		sqlDB.SetConnMaxLifetime(timeout)
	}
	return db, nil
}

// InitRedisDB 初始化redis数据库连接
func InitRedisDB(cfg *config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout * time.Second,
		ReadTimeout:  cfg.ReadTimeout * time.Second,
		WriteTimeout: cfg.WriteTimeout * time.Second,
	})
}

// GetRedis 统一获取redisdb
func GetRedis() *redis.Client {
	//加载配置
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to loa redis config: %v", err)
	}
	//初始化数据库
	db := InitRedisDB(&cfg.Redis)
	return db
}

// GetDB 统一获取mysqldb
func GetDB() *gorm.DB {

	//加载配置
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load mysql config: %v", err)
	}
	//初始化数据库
	db, err := InitDB(&cfg.MySQL)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	return db
}

// Publish 发布消息到redis
func Publish(ctx context.Context, channel string, msg string) error {
	err := GetRedis().Publish(ctx, channel, msg).Err()
	return err
}

// Subscribe 订阅redis消息
func Subscribe(ctx context.Context, channel string) (string, error) {
	sub := GetRedis().Subscribe(ctx, channel)
	for {
		msg, err := sub.ReceiveMessage(ctx)
		log.Printf("Received Redis message: [%s] %s", msg.Channel, msg.Payload)
		if err != nil {
			log.Fatalf("Failed to sub redis message: %v", err)
		}
		return msg.Payload, err
	}
}

const (
	PublishKey = "webSocket"
)

// ListenRedisChannel ai生成的代码（监听频道）
func ListenRedisChannel(ctx context.Context, channel string, wsConn *websocket.Conn) {
	if err := GetRedis().Ping(ctx).Err(); err != nil {
		log.Fatal("Redis 连接失败: ", err)
	}
	sub := GetRedis().Subscribe(ctx, channel)
	// 添加订阅确认日志
	if _, err := sub.Receive(ctx); err != nil {
		log.Printf("订阅失败: %v", err)
		return
	}
	log.Printf("成功订阅频道: %s", channel)
	defer func(sub *redis.PubSub) {
		err := sub.Close()
		if err != nil {
			return
		}
	}(sub)

	/*	ch := sub.Channel()

		log.Printf("Listening on Redis channel: %s", channel)*/

	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Printf("接收消息失败: %v", err)
			return
		}
		log.Printf("Received Redis message: [%s] %s", msg.Channel, msg.Payload)
		// 推送给 WebSocket 客户端
		if err := wsConn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Printf("WebSocket write error: %v", err)
			continue // 不要直接 return，继续监听
		}
	}
}
