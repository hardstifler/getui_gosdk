package getui

import (
	"context"
	"fmt"
	"os"
	"testing"
)

func TestPushClient_PushSingleByCid(t *testing.T) {
	appSecret := "your appSecret"
	appID := "your appID"
	appKey := "your appkey"
	mSecret := "your secret"

	app := NewApp(appID, appKey, appSecret, mSecret)
	p, err := app.MewPushClient()
	if err != nil {
		os.Exit(1)
	}
	ctx := context.TODO()
	requestID := "63549683-a9e-40f1-cu8"
	msg := &PushMessageRequest{
		RequestID: requestID,
		Audience: &Audience{
			Cid: []string{"75ebebdb54b59485c1f45ae2337b7479"},
		},
		PushMessage: &PushMessage{
			Notification: &GeTuiNotification{
				Title:     "哈哈哈",
				Body:      "哈哈",
				ClickType: "url",
				Url:       "www.baidu.com",
			},
		},
	}
	resp, err := p.PushSingleByCid(ctx, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
}
