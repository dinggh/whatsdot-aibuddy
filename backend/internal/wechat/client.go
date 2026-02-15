package wechat

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type Client struct {
	AppID   string
	Secret  string
	APIBase string
	HTTP    *http.Client

	mu         sync.Mutex
	cachedAT   string
	atExpireAt time.Time
}

type code2SessionResp struct {
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
	OpenID     string `json:"openid"`
	SessionKey string `json:"session_key"`
	UnionID    string `json:"unionid"`
}

type Session struct {
	OpenID     string
	SessionKey string
	UnionID    string
}

type accessTokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type phoneResp struct {
	ErrCode   int    `json:"errcode"`
	ErrMsg    string `json:"errmsg"`
	PhoneInfo struct {
		PhoneNumber     string `json:"phoneNumber"`
		PurePhoneNumber string `json:"purePhoneNumber"`
		CountryCode     string `json:"countryCode"`
	} `json:"phone_info"`
}

func New(appID, secret, base string) *Client {
	return &Client{
		AppID:   appID,
		Secret:  secret,
		APIBase: base,
		HTTP: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) Code2Session(ctx context.Context, code string) (Session, error) {
	if code == "" {
		return Session{}, errors.New("empty code")
	}
	u, _ := url.Parse(c.APIBase + "/sns/jscode2session")
	q := u.Query()
	q.Set("appid", c.AppID)
	q.Set("secret", c.Secret)
	q.Set("js_code", code)
	q.Set("grant_type", "authorization_code")
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return Session{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Session{}, err
	}

	var data code2SessionResp
	if err := json.Unmarshal(body, &data); err != nil {
		return Session{}, err
	}
	if data.ErrCode != 0 {
		return Session{}, fmt.Errorf("code2session failed: %d %s", data.ErrCode, data.ErrMsg)
	}
	if data.OpenID == "" {
		return Session{}, errors.New("empty openid")
	}

	return Session{OpenID: data.OpenID, SessionKey: data.SessionKey, UnionID: data.UnionID}, nil
}

func (c *Client) GetPhoneNumberByCode(ctx context.Context, code string) (string, error) {
	if code == "" {
		return "", errors.New("empty code")
	}
	accessToken, err := c.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	u := fmt.Sprintf("%s/wxa/business/getuserphonenumber?access_token=%s", c.APIBase, url.QueryEscape(accessToken))
	payload, _ := json.Marshal(map[string]string{"code": code})
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data phoneResp
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if data.ErrCode != 0 {
		return "", fmt.Errorf("get phone failed: %d %s", data.ErrCode, data.ErrMsg)
	}
	if data.PhoneInfo.PurePhoneNumber != "" {
		return data.PhoneInfo.PurePhoneNumber, nil
	}
	return data.PhoneInfo.PhoneNumber, nil
}

func (c *Client) getAccessToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.cachedAT != "" && time.Now().Before(c.atExpireAt) {
		tk := c.cachedAT
		c.mu.Unlock()
		return tk, nil
	}
	c.mu.Unlock()

	u, _ := url.Parse(c.APIBase + "/cgi-bin/token")
	q := u.Query()
	q.Set("grant_type", "client_credential")
	q.Set("appid", c.AppID)
	q.Set("secret", c.Secret)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	var data accessTokenResp
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}
	if data.ErrCode != 0 {
		return "", fmt.Errorf("get access_token failed: %d %s", data.ErrCode, data.ErrMsg)
	}

	expires := time.Duration(data.ExpiresIn-120) * time.Second
	if expires <= 0 {
		expires = 60 * time.Second
	}

	c.mu.Lock()
	c.cachedAT = data.AccessToken
	c.atExpireAt = time.Now().Add(expires)
	c.mu.Unlock()

	return data.AccessToken, nil
}
