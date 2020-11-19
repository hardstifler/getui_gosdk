package getui

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"getui_gosdk/auth"
	"getui_gosdk/consts"
	"getui_gosdk/httpclient"
	"sync"
	"time"
)

//个推注册App
type App struct {
	AppID        string `json:"app_id"`
	AppKey       string `json:"app_key"`
	AppSecret    string `json:"app_secret"`
	MasterSecret string `json:"master_secret"`
}

//推送客户端
type PushClient struct {
	token string                 //token
	auth  *auth.AuthClient       //获取token
	cli   *httpclient.HTTPClient //发送推送请求
	lock  *sync.RWMutex          //锁
}

func NewApp(appId, appKey, appSecret, masterSecret string) *App {
	return &App{
		AppID:        appId,
		AppKey:       appKey,
		AppSecret:    appSecret,
		MasterSecret: masterSecret,
	}
}

func (a *App) MewPushClient() (*PushClient, error) {
	authCli := auth.NewAuthClient(a.AppID, a.AppKey, a.MasterSecret)
	httpCli := httpclient.NewClient()
	l := &sync.RWMutex{}
	pushCli := &PushClient{
		auth: authCli,
		cli:  httpCli,
		lock: l,
	}
	ctx := context.Background()
	authInfo, err := pushCli.auth.RefreshToken(ctx)
	if err != nil {
		panic(err)
	}
	pushCli.setToken(authInfo.Token)
	fmt.Println(authInfo.Token)
	go func() {
		d := (24*60*60 - 30) * time.Second
		timer := time.NewTimer(d)
		for {
			<-timer.C
			authInfo, err := pushCli.auth.RefreshToken(ctx)
			if err != nil {
				timer.Reset(1 * time.Second)
				continue
			}
			pushCli.setToken(authInfo.Token)
			timer.Reset(d)
		}
	}()
	return pushCli, nil
}

func (p *PushClient) setToken(token string) {
	p.lock.Lock()
	if token != "" {
		p.token = token
	}
	p.lock.Unlock()
}

func (p *PushClient) getToken() string {
	p.lock.RLock()
	token := p.token
	p.lock.RUnlock()
	return token
}

//向单个用户推送消息，可根据cid指定用户
func (p *PushClient) PushSingleByCid(ctx context.Context, request *PushMessageRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/single/cid",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//通过别名推送消息
func (p *PushClient) PushSingleByAlias(ctx context.Context, request *PushMessageRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/single/alias",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//批量发送单推消息，每个cid用户的推送内容都不同的情况下，使用此接口
func (p *PushClient) BatchPushByCid(ctx context.Context, request *BatchPushMessageRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/single/batch/cid",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//批量发送单推消息，在给每个别名用户的推送内容都不同的情况下，可以使用此接口
func (p *PushClient) BatchPushByAlias(ctx context.Context, request *BatchPushMessageRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/single/batch/alias",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//此接口用来创建消息体，并返回taskid，为批量推的前置步骤,taskID需从response解析出来
func (p *PushClient) CreateMessageForBatch(ctx context.Context, request *CreateMessageRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/list/message",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//对列表中所有cid进行消息推送,需要预先创建消息，并带上taskID
func (p *PushClient) BatchPushToCidList(ctx context.Context, request *BatchPushRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/list/cid",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//对列表中所有alias进行消息推送,需要预先创建消息，并带上taskID
func (p *PushClient) BatchPushToAliasList(ctx context.Context, request *BatchPushRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/list/alias",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//群推
func (p *PushClient) PushAll(ctx context.Context, request *PushAllRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/all",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//根据用户tag过滤并推送,Audience.tag不能为空
func (p *PushClient) FilterAndPush(ctx context.Context, request *FilterAndPushRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/tag",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//根据标签[fast_custom_tag]过滤用户并推送。支持定时、定速功能
func (p *PushClient) PushByCustomerTag(ctx context.Context, request *PushByFastTagRequest) (*httpclient.Response, error) {
	if !validateMessage(request) {
		return nil, errors.New("参数不合法")
	}
	data, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	r := &httpclient.Request{
		Method: "POST",
		Path:   p.auth.AppID + "/push/fast_custom_tag",
		Body:   data,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//停止推送任务
func (p *PushClient) StopPushTask(ctx context.Context, taskID string) (*httpclient.Response, error) {
	r := &httpclient.Request{
		Method: "DELETE",
		Path:   p.auth.AppID + "/task/" + taskID,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//查询推送任务
func (p *PushClient) QueryPushTask(ctx context.Context, taskID string) (*httpclient.Response, error) {
	r := &httpclient.Request{
		Method: "GET",
		Path:   p.auth.AppID + "/task/schedule/" + taskID,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

//删除推送任务
func (p *PushClient) DeletePushTask(ctx context.Context, taskID string) (*httpclient.Response, error) {
	r := &httpclient.Request{
		Method: "DELETE",
		Path:   p.auth.AppID + "/task/schedule/" + taskID,
		Header: []httpclient.HTTPOption{
			httpclient.SetHeader("content-type", "application/json;charset=utf-8"),
			httpclient.SetHeader("token", p.getToken()),
		},
	}
	return p.sendRequest(ctx, r)
}

func (p *PushClient) sendRequest(ctx context.Context, request *httpclient.Request) (*httpclient.Response, error) {
	var (
		resp *httpclient.Response
		err  error
	)
	resp, err = p.cli.Do(ctx, request)
	if resp.Code == consts.TokenExpireCode {
		authinfo, err := p.auth.RefreshToken(ctx)
		if err != nil {
			return nil, err
		}
		p.setToken(authinfo.Token)
		resp, err = p.cli.Do(ctx, request)
	}
	return resp, err
}

func validateMessage(msg interface{}) bool {
	if msg == nil {
		return false
	}
	cidMsg, ok := msg.(*PushMessageRequest)
	if ok {
		return len(cidMsg.RequestID) > 0 && cidMsg.Audience != nil &&
			(len(cidMsg.Audience.Cid) > 0 || len(cidMsg.Audience.Alias) > 0 || len(cidMsg.Audience.Tag) > 0 || len(cidMsg.Audience.FastCustomTag) > 0) &&
			cidMsg.PushMessage != nil
	}
	pushAll, ok := msg.(*PushAllRequest)
	if ok {
		return len(pushAll.RequestID) > 0 && pushAll.Audience == "all" && cidMsg.PushMessage != nil
	}

	fp, ok := msg.(*FilterAndPushRequest)
	if ok {
		return len(fp.RequestID) > 0 && fp.Audience != nil && len(fp.Audience.Tag) > 0
	}

	pr, ok := msg.(*PushByFastTagRequest)
	if ok {
		return len(pr.RequestID) > 0 && pr.Audience != nil && len(fp.Audience.FastCustomTag) > 0
	}
	return true
}
