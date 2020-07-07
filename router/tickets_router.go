package router

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"zendesk-integration/handler"
)

func TicketsRouter(e *echo.Echo, client *http.Client) {
	handler := handler.TicketsHandler{
		Client: client,
	}
	e.POST("/tickets", handler.CreateTicket)
}
