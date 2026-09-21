package api

import (
	"bytes"
	"context"
	"encoding/json"
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
	BotTokenHaeder  = "BotToken"
	MaxUserIdHeader = "MaxUserId"
)

func New(client *http.Client, baseUrl string, botToken string) *clientApi {
	return &clientApi{
		httpClient: client,
		baseUrl:    baseUrl,
	}
}

func (c *clientApi) PutUser(ctx context.Context, maxUserID, maxChatID int64) (*domain.User, error) {
	path := "/api/v1/users/me"
	user := User{
		MaxUserId: maxUserID,
		MaxChatId: maxChatID,
	}
	jsonBody, err := json.Marshal(user)
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
	req.Header.Set(MaxUserIdHeader, strconv.Itoa(int(maxUserID)))
	req.Header.Set(BotTokenHaeder, c.botToken)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil
	}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &domain.User{
		ID:        user.ID,
		MaxUserId: user.MaxUserId,
		MaxChatId: user.MaxChatId,
	}, nil
}
