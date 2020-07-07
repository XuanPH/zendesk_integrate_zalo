package impl

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
	"zendesk-integration/config"
	"zendesk-integration/encrypt"
	"zendesk-integration/helper"
	"zendesk-integration/lib"
	"zendesk-integration/model"
	repo "zendesk-integration/respository"

	"github.com/fatih/color"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type IntegrationRepoImpl struct {
	Client *http.Client
	Db     *sqlx.DB
}

func NewIntegrationRepoImpl(client *http.Client, db *sqlx.DB) repo.IntegrationRepository {
	return &IntegrationRepoImpl{
		Client: client,
		Db:     db,
	}
}

func (i *IntegrationRepoImpl) AdminUIHtml(req model.IntegrationRequest, metadata model.Metadata, warning string) string {
	var isWarning string = "false"
	authProfile, err := json.Marshal(metadata.AuthorProfile)

	if err != nil {
		authProfile = []byte("{}")
	}
	if warning != "" {
		isWarning = "true"
	}
	switch {
	case metadata.Priority == "":
		metadata.Priority = "normal"
		fallthrough
	case metadata.Type == "":
		metadata.Type = "question"
	}
	adminUI := helper.ReadAdminUI()
	adminUI = strings.ReplaceAll(adminUI, "#name#", req.Name)
	adminUI = strings.ReplaceAll(adminUI, "#token#", metadata.Token)
	adminUI = strings.ReplaceAll(adminUI, "#returnUrl#", req.ReturnUrl)
	adminUI = strings.ReplaceAll(adminUI, "#priority#", metadata.Priority)
	adminUI = strings.ReplaceAll(adminUI, "#type#", metadata.Type)
	adminUI = strings.ReplaceAll(adminUI, "#tags#", strings.Join(metadata.Tags, ","))
	adminUI = strings.ReplaceAll(adminUI, "#isWarning#", isWarning)
	adminUI = strings.ReplaceAll(adminUI, "#warningMessage#", warning)
	adminUI = strings.ReplaceAll(adminUI, `"#authProfile#"`, string(authProfile))
	adminUI = strings.ReplaceAll(adminUI, `#z_locale#`, string(metadata.Locale))
	adminUI = strings.ReplaceAll(adminUI, `#z_accessToken#`, string(metadata.ZendeskAccessToken))
	adminUI = strings.ReplaceAll(adminUI, `#z_subdomain#`, string(metadata.Subdomain))
	adminUI = strings.ReplaceAll(adminUI, `#z_pushId#`, string(metadata.InstancePushId))
	adminUI = strings.ReplaceAll(adminUI, `#redirect_zalo_url#`, config.ZALO_CALLBACK_URL)
	return adminUI
}

func (i *IntegrationRepoImpl) AdminUISub(metadata model.Metadata, context context.Context) string {
	uri := fmt.Sprintf("%s/getoa?access_token=%s", config.ZALO_END_POINT, metadata.Token)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		adminHtml := i.AdminUIHtml(model.IntegrationRequest{
			Name:      metadata.Name,
			ReturnUrl: metadata.ReturnUrl,
		}, metadata, "Không kết nối được server, Vui lòng kiểm tra lại")
		return adminHtml
	}
	resp, err := i.Client.Do(req)
	if err != nil {
		adminHtml := i.AdminUIHtml(model.IntegrationRequest{
			Name:      metadata.Name,
			ReturnUrl: metadata.ReturnUrl,
		}, metadata, "Không kết nối được server, Vui lòng kiểm tra lại")
		return adminHtml
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		adminHtml := i.AdminUIHtml(model.IntegrationRequest{
			Name:      metadata.Name,
			ReturnUrl: metadata.ReturnUrl,
		}, metadata, "Không kết nối được server, Vui lòng kiểm tra lại")
		return adminHtml
	}
	oaResp := model.OAResponse{}
	if err := json.Unmarshal(body, &oaResp); err != nil || oaResp.Error != 0 {
		adminHtml := i.AdminUIHtml(model.IntegrationRequest{
			Name:      metadata.Name,
			ReturnUrl: metadata.ReturnUrl,
		}, metadata, fmt.Sprintf(`Token %s không đúng, vui lòng kiểm tra lại`, metadata.Token))
		return adminHtml
	}
	new_metadata, _ := json.Marshal(model.Metadata{
		Name:          metadata.Name,
		Token:         metadata.Token,
		Tags:          metadata.Tags,
		Priority:      metadata.Priority,
		Type:          metadata.Type,
		ReturnUrl:     metadata.ReturnUrl,
		Author:        oaResp.Data.OaID,
		AuthorProfile: oaResp.Data,
	})
	htmlRequestForZendesk := fmt.Sprintf(`<html><body>
				<form id="finish"
					method="post"
					action='%s'>
				<input type="hidden"
						name="name"
						value='%s'>
				<input type="hidden"
						name="metadata"
						value='%s'>
				</form>
				<script type="text/javascript">
				var form = document.forms['finish'];
				form.submit();
				</script>
			</body></html>`, metadata.ReturnUrl, metadata.Name, string(new_metadata))
	go i.updateZaloAuthInfo(metadata, oaResp.Data.OaID, metadata.Name)
	return htmlRequestForZendesk
}

func (i *IntegrationRepoImpl) Manifest() interface{} {
	urls := model.ManifestUrls{
		AdminUI:         "./admin_ui",
		PullURL:         "./pull",
		ChannelbackURL:  "./channelback",
		ClickthroughURL: "./clickthrough",
		HealthcheckURL:  "./healthcheck",
	}
	manifest := model.Manifest{
		Name:             "Zalo Channel",
		ID:               "com.quesera.integration.channel.zalo",
		PushClientId:     "zdg-global-quesera-integrations",
		Author:           "Quesera",
		Version:          "v0.0.1",
		ChannelbackFiles: true,
		Urls:             urls,
	}
	return manifest
}

func transformChat(chats []model.DataMessage, metadata model.Metadata) []model.TransformChat {
	oaId := metadata.AuthorProfile.OaID
	dataChat := make([]model.TransformChat, 0)
	for i := range chats {
		chat := chats[i]
		var chatHTML string
		chatHTML = chat.Message
		switch {
		case chat.Type != "text":
			chat.Message = fmt.Sprintf("Send %s", chat.Type)
			fallthrough
		case chat.Type == "photo" || chat.Type == "image" || chat.Type == "gif":
			chatHTML = fmt.Sprintf(`<p>%s</p> <img style='width:315px' src='%s' />`, chat.Message, chat.URL)
			break
		case chat.Type == "voice" || chat.Type == "audio":
			chatHTML = fmt.Sprintf(`<p>%s</p> <a download='audio_file.amr'  href='%s' >Tải xuống</a>`, chat.Message, chat.URL)
			break
		case chat.Type == "locaion":
			position := struct {
				Longitude string `json:"longitude"`
				Latitude  string `json:"latitude"`
			}{
				Latitude:  "",
				Longitude: "",
			}
			json.Unmarshal([]byte(chat.Location), &position)
			chatHTML = fmt.Sprintf(`<p>%s/p> <a target='_blank'  href='http://maps.google.com/maps/place/%s,%s' >Map</a>`, chat.Message, position.Latitude, position.Longitude)
			break
		}
		chatHTML += "" //by pass complier
		fields := []model.TransformChatFields{
			model.TransformChatFields{
				ID:    "type",
				Value: metadata.Type,
			},
			model.TransformChatFields{
				ID:    "priority",
				Value: metadata.Priority,
			},
			model.TransformChatFields{
				ID:    "tags",
				Value: metadata.Tags,
			},
		}

		author := model.TransformChatAuthor{
			ExternalID: strconv.FormatInt(chat.FromID, 10),
			Name:       chat.FromDisplayName,
			ImageURL:   chat.FromAvatar,
		}
		unixTime := unixToDate(chat.Time)
		chatTransform := model.TransformChat{
			ExternalID:  transExternalChatId(strconv.FormatInt(chat.FromID, 10), strconv.FormatInt(chat.ToID, 10), strconv.FormatInt(oaId, 10), chat.MessageID),
			Message:     strip.StripTags(chatHTML),
			HTMLMessage: chatHTML,
			ThreadID:    transExternalChatId(strconv.FormatInt(chat.FromID, 10), strconv.FormatInt(chat.ToID, 10), strconv.FormatInt(oaId, 10), "parent"),
			CreatedAt:   unixTime,
			Author:      author,
			Fields:      fields,
		}
		dataChat = append(dataChat, chatTransform)
	}
	return dataChat
}

func (i *IntegrationRepoImpl) transformPushChat(message model.MessagePush, modelInfo model.ZaloAuthInfo, moreTags []string) interface{} {
	oaId := message.Recipient.ID
	chatHTML := ""
	messageId := message.Message.MsgID
	switch {
	case message.EventName == "user_send_text":
		chatHTML = fmt.Sprintf("<p>%s</p>", message.Message.Text)
		messageId = message.Message.MsgID
		break
	case message.EventName == "user_send_image" || message.EventName == "user_send_gif":
		chatHTML = fmt.Sprintf(`<p style="display:none">Send image</p>`)
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<img style='width:315px' src='%s' />`, attachment.Payload.URL)
		}
		break
	case message.EventName == "user_send_link":
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<a target='_blank'  href='%s' >%s</a>`, attachment.Payload.URL, attachment.Payload.Description)
		}
		break
	case message.EventName == "user_send_audio":
		chatHTML = fmt.Sprintf(`<p style="display:none">Send audio</p>`)
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<a download='audio_file.amr' href='%s' >File Audio</a>`, attachment.Payload.URL)
		}
		break
	case message.EventName == "user_send_video":
		chatHTML = fmt.Sprintf(`<p style="display:none">Send video</p>`)
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<a target="_blank" href="%s">
			<img src="%s" alt="Video" style="width:200px"></a>`, attachment.Payload.URL, attachment.Payload.Thumbnail)
		}
		break
	case message.EventName == "user_send_sticker":
		chatHTML = fmt.Sprintf(`<p style="display:none">Send stickers</p>`)
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<img style='width:100px' src='%s' />`, attachment.Payload.URL)
		}
		break
	case message.EventName == "user_send_location":
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<a target='_blank'  href='http://maps.google.com/maps/place/%s,%s' >Map</a>`, attachment.Payload.Coordinates.Latitude, attachment.Payload.Coordinates.Longitude)
		}
		break
	case message.EventName == "user_send_business_card":
		chatHTML = fmt.Sprintf(`<p>Send business card (not support zendesk type)/p>`)
		break
	case message.EventName == "user_send_file":
		for i := range message.Message.Attachments {
			attachment := message.Message.Attachments[i]
			chatHTML += fmt.Sprintf(`<div style='padding: 10px 0px;'> <a style=" 
			background: rgba(178, 178, 183, 0.36);
			padding: 3px 50px 3px 10px;
			border-radius: 4px;
			font-size: 11px;" target='_blank'  href='%s'>File: %s</a></div>`, attachment.Payload.URL, attachment.Payload.Name)
		}
		break
	}
	var metadata_db model.Metadata
	err := json.Unmarshal([]byte(modelInfo.Metadata), &metadata_db)
	if err != nil {
		metadata_db = model.Metadata{
			Type:     "question",
			Priority: "normal",
			Tags:     []string{"Zalo"},
			Token:    "",
		}
	}
	newTags := metadata_db.Tags
	if moreTags != nil && len(moreTags) > 0 {
		newTags = append(newTags, moreTags...)
	}
	fields := []model.TransformChatFields{
		model.TransformChatFields{
			ID:    "type",
			Value: metadata_db.Type,
		},
		model.TransformChatFields{
			ID:    "priority",
			Value: metadata_db.Priority,
		},
		model.TransformChatFields{
			ID:    "tags",
			Value: newTags,
		},
	}
	authProfile := i.getProfile(metadata_db.Token, message.Sender.ID)
	author := model.TransformChatAuthor{
		ExternalID: message.Sender.ID,
		Name:       authProfile.Data.DisplayName,
		ImageURL:   authProfile.Data.Avatar,
	}
	times, _ := strconv.ParseInt(message.Timestamp, 10, 64)
	unixTime := unixToDate(times)

	chatTransform := model.TransformChat{
		ExternalID:       transExternalChatId(message.Sender.ID, message.Recipient.ID, oaId, messageId),
		Message:          strip.StripTags(chatHTML),
		HTMLMessage:      chatHTML,
		ThreadID:         transExternalChatId(message.Sender.ID, message.Recipient.ID, oaId, "parent"),
		CreatedAt:        unixTime,
		Author:           author,
		Fields:           fields,
		AllowChannelback: true,
	}
	channelPush := struct {
		InstancePushID    string                `json:"instance_push_id"`
		RequestID         string                `json:"request_id"`
		ExternalResources []model.TransformChat `json:"external_resources"`
	}{
		InstancePushID:    modelInfo.InstancePushId,
		RequestID:         encrypt.UUID(),
		ExternalResources: []model.TransformChat{chatTransform},
	}
	return channelPush
}

func transExternalChatId(fromId string, toId string, oaId string, messageId string) string {
	if fromId == oaId {
		return fmt.Sprintf("%s:%s:%s", string(fromId), string(toId), string(messageId))
	}
	return fmt.Sprintf("%s:%s:%s", string(fromId), string(toId), string(messageId))
}

func pullState(chats []model.DataMessage) model.State {
	if len(chats) == 0 || len(chats) < 10 {
		return model.State{
			Offset: 0,
		}
	}
	return model.State{
		Offset: len(chats),
	}
}

func parseExternalChatId(parentId string) []string {
	//0 :oaid
	//1: userid
	//2: messageId
	return strings.Split(parentId, ":")
}

func (i *IntegrationRepoImpl) Pull(metadata model.Metadata, state model.State, c echo.Context) error {
	msgModel, returnCode, err := lib.GetZaloRecentChat(i.Client, metadata.Token, state.Offset)
	if err != nil {
		return c.JSON(returnCode, err.Error())
	}
	chatsAfterTransform := transformChat(msgModel.Data, metadata)
	newState, err := json.Marshal(pullState(msgModel.Data))
	dataResp := struct {
		ExternalResources []model.TransformChat `json:"external_resources"`
		State             string                `json:"state"`
	}{
		ExternalResources: chatsAfterTransform,
		State:             string(newState),
	}
	return c.JSON(http.StatusOK, dataResp)
}

func (i *IntegrationRepoImpl) Push(data model.MessagePush) {
	var modelInfos []model.ZaloAuthInfo
	qrs := fmt.Sprintf("SELECT id, oaId, name, metadata, instancePushId, zendeskAccessToken, subdomain, locale FROM zalo_auth_info where oaId = '%s'", data.Recipient.ID)
	err := i.Db.Select(&modelInfos, qrs)
	if err == nil && len(modelInfos) > 0 {
		FirstMetadata := model.Metadata{}
		json.Unmarshal([]byte(modelInfos[0].Metadata), &FirstMetadata)
		isMessageFromClick, tag := i.checkActionMapping(data.Message.Text)
		if isMessageFromClick == true {
			data.Message.Text = fmt.Sprintf("@TAG:[%s]", strings.Join(tag, ","))
			//send back message to zalo

			recipient := model.SendMessageBodyUser{
				UserID: data.Sender.ID,
			}
			message := model.SendMessageBodyMessage{
				Text: "Cảm ơn bạn đã liên hệ với chúng tôi, bộ phận hỗ trợ sẽ liên lạc với bạn. Cảm ơn",
			}
			postData, _ := json.Marshal(model.SendMessageBody{
				Recipient: recipient,
				Message:   message,
			})
			lib.SendMessageRecipient(i.Client, FirstMetadata.Token, postData)
		}
		isRequestTakeTag := i.isRequestTakeTags(data.Message.Text)
		if isMessageFromClick == false && isRequestTakeTag == true {
			recipient := model.SendMessageBodyUser{
				UserID: data.Sender.ID,
			}
			nilledDefaultAction := model.DefaultAction{
				Type: "oa.open.url",
				URL:  "https://aquavietnam.com.vn/",
			}
			defaultAction1 := model.DefaultAction{
				Type:    "oa.query.hide",
				Payload: "#b3338f2f-1cdb-48a1-ac9d-f0fadeb66826",
			}
			defaultAction2 := model.DefaultAction{
				Type:    "oa.query.hide",
				Payload: "#e80346ff-a0df-48f2-8863-82966ccaa1ac",
			}

			defaultAction3 := model.DefaultAction{
				Type:    "oa.query.hide",
				Payload: "#e5ba6402-249e-4e3d-a24d-472dc4cebe24",
			}
			elements := []model.Elements{
				{Title: "Aqua Support", Subtitle: "Bạn đang cần hỗ trợ gì, vui lòng chọn bên dưới để chúng tôi hỗ trợ bạn tốt hơn", ImageURL: "https://scontent.fsgn2-1.fna.fbcdn.net/v/t1.0-9/80792431_1174120479460840_4539927835149598720_n.jpg?_nc_cat=107&_nc_sid=09cbfe&_nc_ohc=1O-ng39v7l0AX_JVjF3&_nc_ht=scontent.fsgn2-1.fna&oh=0733a363de6cb4d2e7464e55cf7c960f&oe=5F0440C6", DefaultAction: nilledDefaultAction},
				{Title: "Tư vấn mua hàng", Subtitle: "", ImageURL: "https://png.pngtree.com/element_our/sm/20180516/sm_5afbe35fd36cc.jpg", DefaultAction: defaultAction1},
				{Title: "Bảo hành", Subtitle: "", ImageURL: "https://png.pngtree.com/element_our/sm/20180516/sm_5afbe35fd36cc.jpg", DefaultAction: defaultAction2},
				{Title: "Khác", Subtitle: "", ImageURL: "https://png.pngtree.com/element_our/sm/20180516/sm_5afbe35fd36cc.jpg", DefaultAction: defaultAction3},
			}
			payloadMessage := model.Payload{
				TemplateType: "list",
				Elements:     elements,
			}
			attachment := model.Attachment{
				Type:    "template",
				Payload: payloadMessage,
			}

			message := model.SendMessageBodyMessageAttachment{
				Text:       "test",
				Attachment: attachment,
			}
			debugModel := model.SendMessageBodyAttachment{
				Recipient: recipient,
				Message:   message,
			}
			fmt.Println(debugModel)
			postData, _ := json.Marshal(debugModel)
			lib.SendMessageRecipient(i.Client, FirstMetadata.Token, postData)
		}
		data.Message.Text = i.keywordFilter(data.Message.Text)
		for _i := range modelInfos {
			modelInfo := modelInfos[_i]
			transForm := i.transformPushChat(data, modelInfo, tag)
			json_data, err := json.Marshal(transForm)
			if err != nil {
				//log.Errorf("[PUSH-ERROR]- %s", err.Error())
			} else {
				uri := fmt.Sprintf("https://%s.zendesk.com/api/v2/any_channel/push.json", modelInfo.SubDomain)
				req, err := http.NewRequest("POST", uri, bytes.NewBuffer(json_data))
				if err != nil {
					//log.Errorf("[PUSH-ERROR]- %s", err.Error())
				}
				req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", modelInfo.ZendeskAccessToken))
				req.Header.Add("Content-Type", "application/json")
				req.Header.Add("Accept", "application/json")
				resp, err := i.Client.Do(req)
				if err != nil {
					//log.Errorf("[PUSH-ERROR]- %s", err.Error())
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					//log.Errorf("[PUSH-ERROR]- %s", err.Error())
				}
				fmt.Printf("[PUSH-SUCESS]- %s", string(body))
				//log.Infof("[PUSH-SUCESS]- %s", string(body))
				//		fmt.Println(string(json_data))
			}
		}
	} else {
		//log.Errorf("[PUSH-ERROR]- %s", err.Error())
	}
}

func (i *IntegrationRepoImpl) Channelback(metadata model.Metadata, parentId string, channelbackMessage string, channelbackAttachmentUrls []string, c echo.Context) error {
	if len(channelbackAttachmentUrls) > 0 {
		channelbackMessage += "\n\n\n(*)Tệp đính kèm:"
		for i := range channelbackAttachmentUrls {
			channelbackMessage += fmt.Sprintf("\n %d) %s", (i + 1), channelbackAttachmentUrls[i])
		}
	}
	recipient := model.SendMessageBodyUser{
		UserID: parseExternalChatId(parentId)[0],
	}
	message := model.SendMessageBodyMessage{
		Text: channelbackMessage,
	}
	postData, _ := json.Marshal(model.SendMessageBody{
		Recipient: recipient,
		Message:   message,
	})

	dataResp, statusCode, err := lib.SendMessageRecipient(i.Client, metadata.Token, postData)
	if err != nil {
		return c.JSON(statusCode, err.Error())
	}
	return c.JSON(http.StatusOK, struct {
		ExternalID string `json:"external_id"`
	}{
		ExternalID: transExternalChatId(strconv.FormatInt(metadata.AuthorProfile.OaID, 10), recipient.UserID, strconv.FormatInt(metadata.AuthorProfile.OaID, 10), dataResp.Data.MessageID),
	})
}

func unixToDate(unixTimeMs int64) string {
	//1583945608615
	unixTime := unixTimeMs / 1000
	unixTimeUTC := time.Unix(int64(unixTime), 0)          //gives unix time stamp in utc
	unitTimeInRFC3339 := unixTimeUTC.Format(time.RFC3339) // converts utc time to RFC3339 format
	return unitTimeInRFC3339
}

func (i *IntegrationRepoImpl) getProfile(accessToken string, userId string) model.Profile {
	profile := model.Profile{}
	uri := fmt.Sprintf(`%s/getprofile?access_token=%s&data={"user_id":"%s"}`, config.ZALO_END_POINT, accessToken, userId)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return profile
	}
	resp, err := i.Client.Do(req)
	if err != nil {
		return profile
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &profile); err != nil {
		return profile
	}
	return profile
}

func (i *IntegrationRepoImpl) updateZaloAuthInfo(row model.Metadata, oaId int64, name string) {
	uuid := encrypt.UUID()
	json_data, _ := json.Marshal(row)
	query := fmt.Sprintf(`CALL p_zendesk_zalo_mapping_save('%s','%d','%s','%s','%s','%s','%s','%s','%s')`, uuid, oaId, name, string(json_data), row.InstancePushId, row.ZendeskAccessToken, row.Subdomain, row.Locale, "0.0.0.0")
	color.Cyan(query)
	_, err := i.Db.Exec(query)
	if err != nil {
		color.Red("Cant not connect to server - Retry after 5 seconds")
		//log.Errorf("[UpdateZaloAuthInfo-ERROR]- %s", err.Error())
		//time.Sleep(5 * time.Second)
		//go i.updateZaloAuthInfo(row, oaId, name)
	}
}

func (i *IntegrationRepoImpl) checkActionMapping(message string) (isActionClick bool, tag []string) {
	for item := range config.ACTION_MAPPING_CLICK {
		if message == fmt.Sprintf("#%s", item) {
			return true, config.ACTION_MAPPING_CLICK[item]
			break
		}
	}
	return false, nil
}

func (i *IntegrationRepoImpl) isRequestTakeTags(message string) bool {
	removeChar, _ := helper.RemoveUnicodeChar(message)
	isContains := strings.Contains(strings.ToLower(removeChar), "ho tro")
	return isContains
}

func (i *IntegrationRepoImpl) keywordFilter(message string) string {
	//standardizedMessage, _ := helper.RemoveUnicodeChar(strings.ToLower(message))
	//for item := range config.KEYWORK_FILTERING {
	//	if strings.Contains(standardizedMessage, item) {
	//		message = fmt.Sprintf("%s %s", message, config.KEYWORK_FILTERING[item])
	//	}
	//}
	return message
}
