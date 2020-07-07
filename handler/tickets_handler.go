package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"zendesk-integration/helper"
)

type TicketCreateRequest struct {
	Title    string `json:"title"`
	HtmlBody string `json:"html_body"`
}

type RequestTicketObj struct {
	Subject string                  `json:"subject"`
	Comment RequestTicketCommentObj `json:"comment"`
	Tags []string `json:"tags"`
}

type RequestTicketCommentObj struct {
	HTMLBody string `json:"html_body"`
}

type ResponseCreate struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type TicketCreatedResponse struct {
	Ticket struct {
		URL        string      `json:"url"`
		ID         int         `json:"id"`
		ExternalID interface{} `json:"external_id"`
		Via        struct {
			Channel string `json:"channel"`
			Source  struct {
				From struct {
				} `json:"from"`
				To struct {
				} `json:"to"`
				Rel interface{} `json:"rel"`
			} `json:"source"`
		} `json:"via"`
		CreatedAt       time.Time     `json:"created_at"`
		UpdatedAt       time.Time     `json:"updated_at"`
		Type            interface{}   `json:"type"`
		Subject         string        `json:"subject"`
		RawSubject      string        `json:"raw_subject"`
		Description     string        `json:"description"`
		Priority        interface{}   `json:"priority"`
		Status          string        `json:"status"`
		Recipient       interface{}   `json:"recipient"`
		RequesterID     int64         `json:"requester_id"`
		SubmitterID     int64         `json:"submitter_id"`
		AssigneeID      interface{}   `json:"assignee_id"`
		OrganizationID  interface{}   `json:"organization_id"`
		GroupID         interface{}   `json:"group_id"`
		CollaboratorIds []interface{} `json:"collaborator_ids"`
		FollowerIds     []interface{} `json:"follower_ids"`
		EmailCcIds      []interface{} `json:"email_cc_ids"`
		ForumTopicID    interface{}   `json:"forum_topic_id"`
		ProblemID       interface{}   `json:"problem_id"`
		HasIncidents    bool          `json:"has_incidents"`
		IsPublic        bool          `json:"is_public"`
		DueAt           interface{}   `json:"due_at"`
		Tags            []interface{} `json:"tags"`
		CustomFields    []struct {
			ID    int64       `json:"id"`
			Value interface{} `json:"value"`
		} `json:"custom_fields"`
		SatisfactionRating struct {
			Score string `json:"score"`
		} `json:"satisfaction_rating"`
		SharingAgreementIds []interface{} `json:"sharing_agreement_ids"`
		Fields              []struct {
			ID    int64       `json:"id"`
			Value interface{} `json:"value"`
		} `json:"fields"`
		FollowupIds             []interface{} `json:"followup_ids"`
		TicketFormID            int64         `json:"ticket_form_id"`
		BrandID                 int64         `json:"brand_id"`
		SatisfactionProbability interface{}   `json:"satisfaction_probability"`
		AllowChannelback        bool          `json:"allow_channelback"`
		AllowAttachments        bool          `json:"allow_attachments"`
	} `json:"ticket"`
	Audit struct {
		ID        int64     `json:"id"`
		TicketID  int       `json:"ticket_id"`
		CreatedAt time.Time `json:"created_at"`
		AuthorID  int64     `json:"author_id"`
		Metadata  struct {
			System struct {
				Client    string `json:"client"`
				IPAddress string `json:"ip_address"`
				Location  string `json:"location"`
				Latitude  int    `json:"latitude"`
				Longitude int    `json:"longitude"`
			} `json:"system"`
			Custom struct {
			} `json:"custom"`
		} `json:"metadata"`
		Events []struct {
			ID          int64         `json:"id"`
			Type        string        `json:"type"`
			AuthorID    int64         `json:"author_id,omitempty"`
			Body        string        `json:"body,omitempty"`
			HTMLBody    string        `json:"html_body,omitempty"`
			PlainBody   string        `json:"plain_body,omitempty"`
			Public      bool          `json:"public,omitempty"`
			Attachments []interface{} `json:"attachments,omitempty"`
			AuditID     int64         `json:"audit_id,omitempty"`
			Value       string        `json:"value,omitempty"`
			FieldName   string        `json:"field_name,omitempty"`
		} `json:"events"`
		Via struct {
			Channel string `json:"channel"`
			Source  struct {
				From struct {
				} `json:"from"`
				To struct {
				} `json:"to"`
				Rel interface{} `json:"rel"`
			} `json:"source"`
		} `json:"via"`
	} `json:"audit"`
}

type TicketsHandler struct {
	Client *http.Client
}

func (ticket *TicketsHandler) CreateTicket(c echo.Context) error {
	defer c.Request().Body.Close()
	ticketContent := TicketCreateRequest{}
	if err := c.Bind(&ticketContent); err != nil || (strings.Trim(ticketContent.Title, " ") == "" || strings.Trim(ticketContent.HtmlBody, " ") == "") {
		return c.JSON(http.StatusBadRequest, ResponseCreate{
			Message: "DataBindNotFound",
			Data:    "",
		})
	}
	ticketComment := RequestTicketCommentObj{HTMLBody: ticketContent.HtmlBody}

	pureComment, err := helper.RemoveUnicodeChar(ticketComment.HTMLBody)
	if err != nil {
		pureComment = ticketComment.HTMLBody
	}
	var tags = make([]string,0)
	tagMapString := map[string]string {
		"tivi|ti vi": "tivi",
		"tu lanh|tulanh": "tulanh",
		"dien thoai|phone|iphone": "phone",
	}
	for item := range tagMapString {
		arrayOfKeyword := strings.Split(item,"|")
		for keyword := range arrayOfKeyword {
			isContains := strings.Contains(pureComment, string(arrayOfKeyword[keyword]))
			if isContains == true {
				tags = append(tags,tagMapString[item])
				break
			}
		}
	}

	ticketRequest := RequestTicketObj{
		Subject: ticketContent.Title,
		Comment: ticketComment,
		Tags:tags,
	}
	ticketMarshall, err := json.Marshal(struct {
		Ticket RequestTicketObj `json:"ticket"`
	}{
		Ticket: ticketRequest,
	})
	fmt.Println(string(ticketMarshall))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ResponseCreate{
			Message: "Server Error",
			Data:    "",
		})
	}
	req, err := http.NewRequest("POST", "https://pdi-quesera.zendesk.com/api/v2/tickets.json", bytes.NewBuffer(ticketMarshall))
	req.Header.Add("Authorization", "Basic ZHp1bmcubmd1eWVuQHF1ZXNlcmEuc2c6emVuZGVzazk5")
	req.Header.Add("Content-Type", "application/json")
	resp, err := ticket.Client.Do(req)
	body, err := ioutil.ReadAll(resp.Body)

	defer resp.Body.Close()

	if err != nil || resp.StatusCode != 201 {
		return c.JSON(resp.StatusCode, body)
	}
	ticketsResp := TicketCreatedResponse{}
	json.Unmarshal(body, &ticketsResp)
	return c.JSON(http.StatusOK, ResponseCreate{
		Message: "TicketCreated",
		Data:    ticketsResp,
	})
}

