package auth

import (
	"context"
	"fmt"
	"testing"
)

var cli *AuthClient

func init() {
	appID := "your appID"
	appKey := "your appkey"
	mSecret := "your secret"
	cli = NewAuthClient(appID, appKey, mSecret)
}

func TestAuthClient_RefreshToken(t *testing.T) {
	ctx := context.TODO()
	auth, err := cli.RefreshToken(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(auth.ExpireTime)
	fmt.Println(auth.Token)
}

func TestAuthClient_DeleteToken(t *testing.T) {
	token := "your token"
	ctx := context.TODO()
	err := cli.DeleteToken(ctx, token)
	if err != nil {
		panic(err)
	}
}
