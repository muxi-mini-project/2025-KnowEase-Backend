package main

import (
	"KnowEase/controllers"
	"KnowEase/dao"
	"KnowEase/middleware"
	"KnowEase/routes"
	"KnowEase/services"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
	}
}

// @title KnowEase API
// @version 1.0
// @description 小知，你的校园助手
// @host localhost:8080
// @BasePath /api
func main() {
	dsn := os.Getenv("DSN")
	if dsn == "" {
		panic("未设置数据库连接，请检查环境变量")
	}
	db := dao.NewDB(dsn)
	ld := dao.NewLikeDao(db)
	ed := dao.NewEmailDao(db)
	ud := dao.NewUserDao(db)
	pd := dao.NewPostDao(db)
	ls := services.NewLikeService(ld, pd)
	es := services.NewEmailService(ed, ud)
	us := services.NewUserService(ud)
	ts := services.NewTokenService()
	ps := services.NewPostService(pd, ld, ud)
	m := middleware.NewMiddleWare(ts)
	lc := controllers.NewLikeControllers(ls, ps, us)
	pc := controllers.NewPostControllers(ps, ls, es, us)
	uc := controllers.NewUserControllers(us, es, ts)
	upsv := routes.NewUserPageSvc(lc, m, pc)
	psv := routes.NewPostSvc(lc, m, pc)
	usv := routes.NewUserSvc(uc, m)
	app := routes.NewApp(usv, psv, upsv)
	app.Run()
	go lc.UpdateAllCount()
}
