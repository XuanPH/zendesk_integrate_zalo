package respository

import (
	"context"
	"zendesk-integration/model"

	"github.com/labstack/echo/v4"
)

type IntegrationRepository interface {
	AdminUIHtml(req model.IntegrationRequest, metadata model.Metadata, warning string) string
	Manifest() interface{}
	AdminUISub(metadata model.Metadata, context context.Context) string
	Pull(metadata model.Metadata, state model.State, c echo.Context) error
	Channelback(metadata model.Metadata, parentId string, channelbackMessage string, channelbackAttachmentUrls []string, c echo.Context) error
	Push(data model.MessagePush)
}
