package router

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"net/http"
	"zendesk-integration/handler"
	"zendesk-integration/respository/impl"
)

func IntegrationRouter(e *echo.Echo, client *http.Client, db *sqlx.DB) {
	handler := handler.IntegrationHandler{
		IntegrationRepo: impl.NewIntegrationRepoImpl(client, db),
	}
	e.GET("/manifest", handler.Manifest)
	e.POST("/admin_ui", handler.AdminUI)
	e.POST("/admin_ui_2", handler.AdminUISub)
	e.POST("/pull", handler.Pull)
	e.POST("/channelback", handler.Channelback)
	e.POST("/push", handler.Push)
	e.GET("/GetToken", handler.GetZaloToken)
}
