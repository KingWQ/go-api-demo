package bootstrap

import (
	"errors"
	"fmt"
	"go-api-demo/app/models/user"
	"go-api-demo/pkg/config"
	"go-api-demo/pkg/database"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

// SetupDB 初始化数据库和 ORM
func SetupDB() {
	var dbConfig gorm.Dialector
	switch config.Get("database.connection") {
	case "mysql":
		//构建 DSN 信息
		dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=%v&parseTime=True&multiStatements=true&loc=Local",
			config.Get("database.mysql.username"),
			config.Get("database.mysql.password"),
			config.Get("database.mysql.host"),
			config.Get("database.mysql.port"),
			config.Get("database.mysql.database"),
			config.Get("database.mysql.charset"),
		)
		dbConfig = mysql.New(mysql.Config{
			DSN: dsn,
		})
	case "sqlite":
		//初始化 sqlite
		sqliteDatabase := config.Get("database.sqlite.database")
		dbConfig = sqlite.Open(sqliteDatabase)
	default:
		panic(errors.New("database connection not supported"))
	}

	//连接数据库，并设置GORM的日志模式
	database.Connect(dbConfig, logger.Default.LogMode(logger.Info))

	//设置最大连接数
	database.SQLDB.SetMaxOpenConns(config.GetInt("database.mysql.max_open_connections"))
	//设置最大空闲连接数
	database.SQLDB.SetMaxIdleConns(config.GetInt("database.mysql.max_idle_connections"))

	//设置每个连接的过期时间
	database.SQLDB.SetConnMaxLifetime(time.Duration(config.GetInt("database.mysql.max_life_seconds")))

	err := database.DB.AutoMigrate(&user.User{})
	if err != nil {
		fmt.Println(err.Error())
	}
}
