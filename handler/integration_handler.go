package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
	"zendesk-integration/helper"
	"zendesk-integration/model"
	"zendesk-integration/respository"

	"github.com/labstack/echo/v4"
)

type IntegrationHandler struct {
	IntegrationRepo respository.IntegrationRepository
}

func (i *IntegrationHandler) Manifest(c echo.Context) error {
	manifest := i.IntegrationRepo.Manifest()
	//log.Info("Someone call manifest")
	return c.JSON(http.StatusOK, manifest)
}

/*
// admin ui post example
{
  "name":"",
  "metadata":"",
  "state":"",
  "return_url":"https://zdg-pdi-quesera.zendesk.com/zendesk/channels/integration_service_instances/editor_finalizer",
  "instance_push_id":"5967536b-f983-449f-8848-05aeb3b2746d",
  "zendesk_access_token":"ea4642ddf1332328f233af81ef77d02bc405399f72d48b4b228e7ad0b89ac4de",
  "subdomain":"zdg-pdi-quesera",
  "locale":"en-US"
}
*/

func (i *IntegrationHandler) AdminUI(c echo.Context) error {
	req := model.IntegrationRequest{}
	metadata := model.Metadata{}
	defer c.Request().Body.Close()
	reqMetadata := c.FormValue("metadata")
	json.Unmarshal([]byte(reqMetadata), &metadata)
	c.Bind(&req)
	metadata.InstancePushId = c.FormValue("instance_push_id")
	metadata.Locale = c.FormValue("locale")
	metadata.ZendeskAccessToken = c.FormValue("zendesk_access_token")
	metadata.Subdomain = c.FormValue("subdomain")
	var AdminHtml = i.IntegrationRepo.AdminUIHtml(req, metadata, "")
	return c.HTML(http.StatusOK, AdminHtml)
}

func (i *IntegrationHandler) AdminUISub(c echo.Context) error {
	metadata := model.Metadata{}
	defer c.Request().Body.Close()
	data := c.Request()
	fmt.Println(data)
	if err := c.Bind(&metadata); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	ctx, _ := context.WithTimeout(c.Request().Context(), 10*time.Second)
	html := i.IntegrationRepo.AdminUISub(metadata, ctx)
	return c.HTML(http.StatusOK, html)
}

func (i *IntegrationHandler) Pull(c echo.Context) error {
	//log.Infof("[POST]- Pull method called")
	req := model.IntegrationRequest{}
	rmetadata := c.FormValue("metadata")
	rstate := c.FormValue("state")
	metadata := model.Metadata{}
	state := model.State{}
	json.Unmarshal([]byte(rmetadata), &metadata)
	json.Unmarshal([]byte(rstate), &state)

	defer c.Request().Body.Close()

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}

	return i.IntegrationRepo.Pull(metadata, state, c)
}

func (i *IntegrationHandler) Channelback(c echo.Context) error {
	//log.Infof("[POST]- Channelback method called")
	defer c.Request().Body.Close()

	type TempRequest struct {
		Metadata                string   `json:"metadata" form:"metadata"`
		ParentId                string   `json:"parent_id" form:"parent_id"`
		Message                 string   `json:"message" form:"message"`
		FileUrls                []string `json:"file_urls" form:"file_urls[]"`
		ThreadID                string   `json:"thread_id" form:"thread_id"`
		RequestUniqueIdentifier string   `json:"request_unique_identifier" form:"request_unique_identifier"`
	}
	req := TempRequest{}
	metadata := model.Metadata{}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusOK, err.Error())
	}
	json.Unmarshal([]byte(req.Metadata), &metadata)
	return i.IntegrationRepo.Channelback(metadata, req.ThreadID, req.Message, req.FileUrls, c)
}

func (i *IntegrationHandler) Push(c echo.Context) error {
	defer c.Request().Body.Close()
	message := model.MessagePush{}
	if err := c.Bind(&message); err != nil {
		fmt.Println(err.Error())
		return c.JSON(http.StatusOK, "")
	}
	i.IntegrationRepo.Push(message)
	return c.JSON(http.StatusOK, "Take it")
}

func (i *IntegrationHandler) GetZaloToken(c echo.Context) error {
	tokenResp := c.QueryParam("access_token")
	//log.Infof("[Get]- GetZaloToken method called")
	oaId := c.QueryParam("oaId")
	html := helper.TokenPage()
	html = strings.ReplaceAll(html, `#token#`, tokenResp)
	html = strings.ReplaceAll(html, `#oaid#`, oaId)
	return c.HTML(http.StatusOK, html)
}
