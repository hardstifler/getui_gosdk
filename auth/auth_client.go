package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"getui_gosdk/httpclient"
	"time"
)

type AuthRequest struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	AppKey    string `json:"appkey"`
}

type AuthInfo struct {
	ExpireTime string `json:"expire_time"`
	Token      string `json:"token"`
}

type AuthClient struct {
	appKey string
	AppID  string
	secret string
	cli    *httpclient.HTTPClient
}

func NewAuthClient(appID string, appKey string, masterSecret string) *AuthClient {
	if len(appKey) <= 0 || len(appID) <= 0 || len(masterSecret) <= 0 {
		panic("appKey或appID或MasterSecret不合法")
	}
	c := httpclient.NewClient()
	return &AuthClient{
		appKey: appKey,
		AppID:  appID,
		secret: masterSecret,
		cli:    c,
	}
}

//刷新token，服务初始化调用，如果接口返回10001是被动调用
func (c *AuthClient) RefreshToken(ctx context.Context) (*AuthInfo, error) {
	timestamp := fmt.Sprintf("%d", time.Now().Unix()*1000)
	sign := genSign(c.appKey, timestamp, c.secret)
	req := &AuthRequest{
		Sign:      sign,
		Timestamp: timestamp,
		AppKey:    c.appKey,
	}
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   c.AppID + "/" + "auth",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
		},
	}
	resp, err := c.cli.Do(ctx, r)
	if err != nil {
		return nil, err
	}
	if resp.Data != nil {
		var authInfo AuthInfo
		if err := json.Unmarshal(resp.Data, &authInfo); err != nil {
			return nil, err
		}
		return &authInfo, nil
	}
	return nil, errors.New("未知错误")
}

//删除鉴权token
func (c *AuthClient) DeleteToken(ctx context.Context, token string) error {
	r := &httpclient.Request{
		Method: "DELETE",
		Path:   c.AppID + "/auth" + "/" + token,
	}
	resp, err := c.cli.Do(ctx, r)
	if err != nil {
		return err
	}
	if resp != nil && resp.Code != 0 {
		return errors.New("未知错误")
	}
	return nil
}

func genSign(appkey string, timeStamp string, masterSecret string) string {
	key := fmt.Sprintf("%s%s%s", appkey, timeStamp, masterSecret)
	h := sha256.New()
	h.Write([]byte(key))
	sum := h.Sum(nil)
	return hex.EncodeToString(sum)
}
