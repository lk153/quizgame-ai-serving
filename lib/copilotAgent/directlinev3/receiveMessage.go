package directlinev3

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func (h *directLine) ReceiveMessages(
	ctx context.Context, conversationID string, watermarkNum int,
) (data ReceiveMessages, err error) {
	watermarkStr := ""
	if watermarkNum > 0 {
		watermarkStr = strconv.Itoa(watermarkNum)
	}
	url := fmt.Sprintf("%s/conversations/%s/activities?watermark=%s",
		strings.TrimRight(ApiPath, "/"), conversationID, watermarkStr)
	req, err := http.NewRequest(HTTP_GET, url, nil)
	if err != nil {
		return
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
