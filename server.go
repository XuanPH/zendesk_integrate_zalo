package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"os"
	"zendesk-integration/config"
	"zendesk-integration/db"
	"zendesk-integration/helper"
	"zendesk-integration/router"
)

func init() {
	os.Setenv("APP_NAME", "zendesk_integration")
	//fmt.Println(">>>>>", os.Getenv("APP_NAME"))
	//log.InitLogger(false)
}

func main() {
	helper.ShowBanner()

	sql := &db.Sql{
		Host:   config.DB_HOST,
		Port:   config.DB_PORT,
		Uid:    config.DB_UID,
		Pwd:    config.DB_PWD,
		DbName: config.DB_NAME,
	}
	sql.Connect()
	defer sql.Close()
	e := echo.New()
	e.HideBanner = true
	DefaultCORSConfig := middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}
	e.Use(middleware.CORSWithConfig(DefaultCORSConfig))
	e.GET("/healthcheck", func(context echo.Context) error {
		return context.String(http.StatusOK, "v0.1.7 - Hide filter keyword - tag v0.0.26")
	})
	e.GET("/logtail", func(context echo.Context) error {
		file := context.QueryParam("file")
		_type := context.QueryParam("type")
		return context.JSON(http.StatusOK, helper.ReadLogFile(_type, file))
	})
	client := &http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       0,
	}
	router.IntegrationRouter(e, client, sql.Db)
	router.TicketsRouter(e, client)
	e.Logger.Fatal(e.Start(":8080"))
}
