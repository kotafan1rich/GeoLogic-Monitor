package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kotafan1rich/GeoLogic-Monitor/bot/internal/domain"
)

type clientApi struct {
	httpClient *http.Client
	baseUrl    string
	botToken   string
}

const (
	AuthorizationHeader = "Authorization"
	MaxUserIDHeader     = "X-Max-User-Id"
)

func New(client *http.Client, baseUrl string, botToken string) *clientApi {
	return &clientApi{
		httpClient: client,
		baseUrl:    baseUrl,
		botToken:   botToken,
	}
}

func (c *clientApi) PutUser(ctx context.Context, maxUserID, maxChatID int64) (*domain.User, error) {
	path := "/api/v1/users/me"
	request := UpsertUserRequest{
		MaxChatId: maxChatID,
	}
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(
		ctx,
		"PUT",
		c.baseUrl+path,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set(MaxUserIDHeader, strconv.FormatInt(maxUserID, 10))
	req.Header.Set(AuthorizationHeader, "Bearer "+c.botToken)
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upsert user: unexpected status %s", resp.Status)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		MaxUserId: user.MaxUserId,
		MaxChatId: user.MaxChatId,
	}, nil
}
