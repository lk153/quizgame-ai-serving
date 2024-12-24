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

func (h *directLine) GenerateToken(ctx context.Context) (data GenerateTokenResp, err error) {
	url := fmt.Sprintf("%s/%s", strings.TrimRight(ApiPath, "/"), "tokens/generate")
	payload := []byte(`{}`)
	req, err := http.NewRequest(HTTP_POST, url, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", secret))
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
	if err == nil {
		h.Token = data.Token
	}

	return
}
