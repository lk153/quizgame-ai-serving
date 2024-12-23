package copilotagent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (h *directLineV3) StartConversation(ctx context.Context) (data CreatedConversation, err error) {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(ApiPath, "/"), "conversations")
	payload := []byte(`{}`)
	req, err := http.NewRequest(HTTP_POST, url, bytes.NewBuffer(payload))
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
