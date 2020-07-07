package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"zendesk-integration/config"
	"zendesk-integration/model"
)

func GetZaloRecentChat(client *http.Client, accessToken string, offSet int) (model.MesageResponse, int, error) {
	uri := fmt.Sprintf(`%s/listrecentchat?access_token=%s&data={"offset":%d,"count":10}`, config.ZALO_END_POINT, accessToken, offSet)
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return model.MesageResponse{}, http.StatusBadRequest, err
	}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return model.MesageResponse{}, http.StatusBadRequest, errDo
	}
	defer resp.Body.Close()
	body, errRead := ioutil.ReadAll(resp.Body)
	if errRead != nil || resp.StatusCode != 200 {
		return model.MesageResponse{}, http.StatusBadRequest, errRead
	}
	msgModel := model.MesageResponse{}
	if err := json.Unmarshal(body, &msgModel); err != nil || msgModel.Error != 0 {
		return model.MesageResponse{}, http.StatusInternalServerError, err
	}
	return msgModel, http.StatusOK, nil
}
func SendMessageRecipient(client *http.Client, accessToken string, postData []byte) (model.SendMessageResponse, int, error) {
	uri := fmt.Sprintf(`%s/message?access_token=%s`, config.ZALO_END_POINT, accessToken)
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(postData))
	if err != nil {
		return model.SendMessageResponse{}, http.StatusBadRequest, err
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return model.SendMessageResponse{}, http.StatusBadRequest, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || resp.StatusCode != 200 {
		return model.SendMessageResponse{}, http.StatusBadRequest, err
	}
	dataResp := model.SendMessageResponse{}
	if err := json.Unmarshal(body, &dataResp); err != nil || dataResp.Error != 0 {
		var errMsg string
		if err != nil {
			errMsg = err.Error()
		} else {
			errMsg = dataResp.Message
		}
		return model.SendMessageResponse{}, http.StatusInternalServerError, errors.New(errMsg)
	}
	return dataResp, http.StatusOK, nil

}
