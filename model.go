package getui

//推送消息体定义

//推送目标，暂时不支持【all】所有条件只能选一个
type Audience struct {
	Cid           []string `json:"cid,omitempty"`             //按照用户clientID推送
	Alias         []string `json:"alias,omitempty"`           //按照用户别名推送
	FastCustomTag string   `json:"fast_custom_tag,omitempty"` //按照标签别名推送
	Tag           []*Tag   `json:"fast_custom_tag,omitempty"` //按照组合条件别名推送
}

//组和条件
type Tag struct {
	Key    string   `json:"key"`
	Values []string `json:"values"`
	OpType string   `json:"opt_type"`
}

//推送设置
type PushSetting struct {
	Ttl          int64          `json:"ttl,omitempty"`           //消息离线时间
	Strategy     map[string]int `json:"strategy,omitempty"`      //厂商通道策略
	Speed        int            `json:"speed,omitempty"`         //控制个推定速推送
	ScheduleTime int            `json:"schedule_time,omitempty"` //定时推送 毫秒 vip功能，未开通不要填写
}

//个推通道消息
type PushMessage struct {
	Duration     string             `json:"duration,omitempty"` //手机端通知展示时间段，格式为毫秒时间戳段，两个时间的时间差必须大于10分钟，例如："1590547347000-1590633747000"
	Notification *GeTuiNotification `json:"notification"`       //安卓通知消息内容，仅支持安卓系统，iOS系统不展示个推通知消息，与transmission二选一，两个都填写时报错
	Transmission string             `json:"transmission"`       //纯透传消息内容，安卓和iOS均支持，与notification 二选一，两个都填写时报错，长度 ≤ 3072
}

//通知消息内容，android only
type GeTuiNotification struct {
	Title        string `json:"title"`                //通知消息标题，长度 ≤ 50
	Body         string `json:"body"`                 //通知消息内容，长度 ≤ 256
	BigText      string `json:"big_text"`             //长文本消息内容，通知消息+长文本样式，与big_image二选一，两个都填写时报错，长度 ≤ 512
	BigImage     string `json:"big_image"`            //大图的URL地址，通知消息+大图样式， 与big_text二选一，两个都填写时报错，长度 ≤ 1024
	Logo         string `json:"logo"`                 //通知的图标名称，包含后缀名（需要在客户端开发时嵌入），如“push.png”，长度 ≤ 64
	LogoUrl      string `json:"logo_url"`             //通知图标URL地址，长度 ≤ 256
	ChannelId    string `json:"channel_id"`           //通知渠道id，长度 ≤ 64
	ChannelName  string `json:"channel_name"`         //通知渠道名称，长度 ≤ 64
	ChannelLevel int    `json:"channel_level"`        //http://docs.getui.com/getui/server/rest_v2/common_args/?id=doc-title-5
	ClickType    string `json:"click_type,omitempty"` //点击通知后续动作[intent,url,payload,startapp,none]
	Intent       string `json:"intent,omitempty"`     //click_type为intent时必填,点击通知打开应用特定页面，长度 ≤ 2048;
	Url          string `json:"url,omitempty"`        //click_type为url时必填,点击通知打开链接，长度 ≤ 1024
	Payload      string `json:"payload,omitempty"`    //click_type为payload时必填,长度 ≤ 3072
	NotifyId     int32  `json:"notify_id,omitempty"`  //覆盖任务时会使用到该字段，两条消息的notify_id相同，新的消息会覆盖老的消息，长度32位
}

//个推厂商通道消息
type PushChannel struct {
	Ios     *IosChannelMessgae     `json:"ios"`     //ios通道推送消息内容
	Android *AndroidChannelMessgae `json:"android"` //android通道推送消息内容
}

//ios厂商通道消息
type IosChannelMessgae struct {
	Type           string        `json:"type"`                       //voip：voip语音推送，notify：apns通知消息 默认notify
	IosAps         Aps           `json:"aps,omitempty"`              //推送通知消息内容
	AutoBadge      string        `json:"auto_badge,omitempty"`       //用于计算icon上显示的数字，还可以实现显示数字的自动增减，如“+1”、 “-1”、 “1” 等，计算结果将覆盖badge
	Payload        string        `json:"payload,omitempty"`          //增加自定义的数据
	Multimedia     IosMultimedia `json:"multimedia"`                 //多媒体设置
	ApnsCollapseId string        `json:"apns-collapse-id,omitempty"` //使用相同的apns-collapse-id可以覆盖之前的消息
}

//ios推送通知消息内容
type Aps struct {
	IosAlert         Alert  `json:"alert"`
	ContentAvailable int    `json:"content-available"` //0 普通； 1 静默
	Sound            string `json:"sound"`             //通知铃声文件名，如果铃声文件未找到，响铃为系统默认铃声。 无声设置为“com.gexin.ios.silence”或不填
	Category         string `json:"category"`          //在客户端通知栏触发特定的action和button显示
	ThreadId         string `json:"thread-id"`         //ios的远程通知通过该属性对通知进行分组，仅支持iOS 12.0以上版本
}

//ios通知消息
type Alert struct {
	Title           string   `json:"title,omitempty"`             //通知消息标题
	Body            string   `json:"body,omitempty"`              //通知消息内容
	ActionLocKey    string   `json:"action-loc-key,omitempty"`    //（用于多语言支持）指定执行按钮所使用的Localizable.strings
	LocKey          string   `json:"loc-key,omitempty"`           //（用于多语言支持）指定Localizable.strings文件中相应的key
	LocArgs         []string `json:"loc-args,omitempty"`          //如果loc-key中使用了占位符，则在loc-args中指定各参数
	LaunchImage     string   `json:"launch-image,omitempty"`      //指定启动界面图片名
	TitleLocKey     string   `json:"title-loc-key,omitempty"`     //(用于多语言支持）对于标题指定执行按钮所使用的Localizable.strings,仅支持iOS8.2以上版本
	TitleLocArgs    []string `json:"title-loc-args,omitempty"`    //对于标题,如果loc-key中使用的占位符，则在loc-args中指定各参数,仅支持iOS8.2以上版本
	Subtitle        string   `json:"subtitle,omitempty"`          //通知子标题,仅支持iOS8.2以上版本
	SubtitleLocKey  string   `json:"subtitle-loc-key,omitempty"`  //当前本地化文件中的子标题字符串的关键字,仅支持iOS8.2以上版本
	SubtitleLocArgs []string `json:"subtitle-loc-args,omitempty"` //当前本地化子标题内容中需要置换的变量参数 ,仅支持iOS8.2以上版本
}

//ios多媒体设置
type IosMultimedia struct {
	Url      string `json:"url"`       //多媒体资源地址
	Type     int    `json:"type"`      //资源类型（1.图片，2.音频，3.视频）
	OnlyWifi bool   `json:"only_wifi"` //是否只在wifi环境下加载，如果设置成true,但未使用wifi时，会展示成普通通知
}

//android厂商通道消息
type AndroidChannelMessgae struct {
	AndroidUps Ups `json:"ups"`
}

//android厂商通道推送消息内容
type Ups struct {
	Notification ChannelNotification //厂商通知消息内容
	Transmission string              `json:"transmission"` //透传消息内容，与notification 二选一，两个都填写时报错，长度 ≤ 3072
}

//厂商通知消息内容
type ChannelNotification struct {
	Title     string           `json:"title"`                //通知消息标题，长度 ≤ 50
	Body      string           `json:"body"`                 //通知消息内容，长度 ≤ 256
	ClickType string           `json:"click_type,omitempty"` //点击通知后续动作[intent,url,payload,startapp,none]
	Intent    string           `json:"intent,omitempty"`     //click_type为intent时必填,点击通知打开应用特定页面，长度 ≤ 2048;
	Url       string           `json:"url,omitempty"`        //click_type为url时必填,点击通知打开链接，长度 ≤ 1024
	Payload   string           `json:"payload,omitempty"`    //click_type为payload时必填,长度 ≤ 3072
	NotifyId  int32            `json:"notify_id,omitempty"`  //覆盖任务时会使用到该字段，两条消息的notify_id相同，新的消息会覆盖老的消息，长度32位
	Options   []*ChannelOption `json:"options"`
}

//第三方厂商通知扩展内容
type ChannelOption struct {
	Constraint string      `json:"constraint"` //扩展内容对应厂商通道设置如：HW,MZ,...
	Key        string      `json:"key"`        //厂商内容扩展字段,单个厂商特有字段，http://docs.getui.com/getui/server/rest_v2/common_args/?id=doc-title-5
	Value      interface{} `json:"value"`      //value的设置根据key值决定。例如，
}

//单推结构体
type PushMessageRequest struct {
	RequestID   string       `json:"request_id"`             //请求唯一标识号，10-32位之间；如果request_id重复，会导致消息丢失
	Audience    *Audience    `json:"audience"`               //推送目标用户，详细解释见下方audience说明
	Settings    *PushSetting `json:"settings,omitempty"`     //推送条件设置，详细解释见下方settings说明
	PushMessage *PushMessage `json:"push_message"`           //个推推送消息参数
	PushChannel *PushChannel `json:"push_channel,omitempty"` //厂商推送消息参数，包含ios消息参数，android厂商消息参数
}

//批量推送结构体;每个cid对应消息不同时使用
type BatchPushMessageRequest struct {
	IsAsync bool                  `json:"is_async"`
	MsgList []*PushMessageRequest `json:"msg_list"`
}

//创建消息结构体
type CreateMessageRequest struct {
	RequestID   string       `json:"request_id"`
	GroupName   string       `json:"group_name"`
	Settings    *PushSetting `json:"settings,omitempty"`     //推送条件设置，详细解释见下方settings说明
	PushMessage *PushMessage `json:"push_message"`           //个推推送消息参数
	PushChannel *PushChannel `json:"push_channel,omitempty"` //厂商推送消息参数，包含ios消息参数，android厂商消息参数
}

//根据cid[alias]群推
type BatchPushRequest struct {
	IsAsync  bool      `json:"is_async"`
	Audience *Audience `json:"audience"` //推送目标用户，详细解释见下方audience说明
	TaskId   string    `json:"task_id"`  //CreateMessageRequest创建消息返回的taskID
}

//全推结构体
type PushAllRequest struct {
	RequestID   string       `json:"request_id"` //请求唯一标识号，10-32位之间；如果request_id重复，会导致消息丢失
	GroupName   string       `json:"group_name"`
	Settings    *PushSetting `json:"settings,omitempty"`     //推送条件设置，详细解释见下方settings说明
	Audience    string       `json:"audience"`               //推送目标用户，详细解释见下方audience说明all表示推全部
	PushMessage *PushMessage `json:"push_message"`           //个推推送消息参数
	PushChannel *PushChannel `json:"push_channel,omitempty"` //厂商推送消息参数，包含ios消息参数，android厂商消息参数
}

//按照tag过滤用户推送
type FilterAndPushRequest struct {
	RequestID   string       `json:"request_id"` //请求唯一标识号，10-32位之间；如果request_id重复，会导致消息丢失
	GroupName   string       `json:"group_name"`
	Settings    *PushSetting `json:"settings,omitempty"`     //推送条件设置，详细解释见下方settings说明
	Audience    *Audience    `json:"audience"`               //推送目标用户，详细解释见下方audience说明all表示推全部
	PushMessage *PushMessage `json:"push_message"`           //个推推送消息参数
	PushChannel *PushChannel `json:"push_channel,omitempty"` //厂商推送消息参数，包含ios消息参数，android厂商消息参数
}

//按照fast_custom_tag推送
type PushByFastTagRequest struct {
	*FilterAndPushRequest
}
