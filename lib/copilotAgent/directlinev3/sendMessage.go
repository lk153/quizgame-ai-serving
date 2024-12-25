package directlinev3

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h *directLine) SendMessage(
	ctx context.Context, conversationID string, userID string, message string,
) (data SendMessageResp, err error) {
	url := fmt.Sprintf("%s/conversations/%s/activities", strings.TrimRight(ApiPath, "/"), conversationID)
	reqPayload := SendMessageReq{
		Locale: DefaultLocale,
		Type:   DefaultMessageType,
		From: From{
			ID:   userID,
			Name: "User",
			Role: "user",
		},
		Text: message,
	}

	payload, err := json.Marshal(reqPayload)
	if err != nil {
		return
	}

	req, err := http.NewRequest(HTTP_POST, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	if strings.EqualFold(h.Token, "") {
		h.Token = secret
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", h.Token))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Transport: http.DefaultTransport}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &data)
	return
}
