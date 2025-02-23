package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

func ProvideDBConnectionString() string {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("未设置数据库连接，请检查环境变量")
	}

	return dsn
}

// @title KnowEase API
// @version 1.0
// @description 小知，你的校园助手
// @host localhost:8080
// @BasePath /api
func main() {
	app := InitializeApp()
	app.Run()
}
