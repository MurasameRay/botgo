package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/tencent-connect/botgo"
	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/dto/message"
	"github.com/tencent-connect/botgo/event"
	"github.com/tencent-connect/botgo/interaction/webhook"
	"github.com/tencent-connect/botgo/token"
	"gopkg.in/yaml.v3"
)

func GetValidFile(w http.ResponseWriter, r *http.Request) {
	// 文件路径
	filePath := "102457514.json" // 替换为你要下载的文件的实际路径

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()
	// 获取文件信息
	fileInfo, err := file.Stat()
	if err != nil {
		http.Error(w, "Could not get file info", http.StatusInternalServerError)
		return
	}

	// 设置响应头
	w.Header().Set("Content-Disposition", "attachment; filename="+fileInfo.Name())
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	// 写入文件内容到响应
	http.ServeContent(w, r, fileInfo.Name(), fileInfo.ModTime(), file)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// 从查询参数中获取 "hello" 的值
	helloValue := r.URL.Query().Get("hello")
	if helloValue == "" {
		http.Error(w, "Missing 'hello' parameter", http.StatusBadRequest)
		return
	}

	// 返回该值
	fmt.Fprintf(w, "Hello, %s!", helloValue)
}

const (
	host_ = "0.0.0.0"
	port_ = 9000
	path_ = "/qqbot"
)

// 消息处理器，持有 openapi 对象
var processor Processor

func main() {
	ctx := context.Background()
	// 加载 appid 和 token
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln("load config file failed, err:", err)
	}
	credentials := &token.QQBotCredentials{
		AppID:     "",
		AppSecret: "",
	}
	if err = yaml.Unmarshal(content, &credentials); err != nil {
		log.Fatalln("parse config failed, err:", err)
	}
	log.Println("credentials:", credentials)
	tokenSource := token.NewQQBotTokenSource(&token.QQBotCredentials{
		AppID:     credentials.AppID,
		AppSecret: credentials.AppSecret,
	})
	if err = token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
		log.Fatalln(err)
	}
	// 初始化 openapi，正式环境
	api := botgo.NewOpenAPI(credentials.AppID, tokenSource).WithTimeout(5 * time.Second).SetDebug(true)
	processor = Processor{api: api}
	// 注册处理函数
	_ = event.RegisterHandlers(
		// ***********消息事件***********
		// 群@机器人消息事件
		GroupATMessageEventHandler(),
		// C2C消息事件
		C2CMessageEventHandler(),
		// 频道@机器人事件
		ChannelATMessageEventHandler(),
	)
	// 注册新的接口
	http.HandleFunc("/hello", helloHandler) // 这里是添加的接口
	http.HandleFunc("/102457514.json", GetValidFile)
	http.HandleFunc(path_, func(writer http.ResponseWriter, request *http.Request) {
		webhook.HTTPHandler(writer, request, credentials)
	})
	if err = http.ListenAndServe(fmt.Sprintf("%s:%d", host_, port_), nil); err != nil {
		log.Fatal("setup server fatal:", err)
	}
}

// ChannelATMessageEventHandler 实现处理 at 消息的回调
func ChannelATMessageEventHandler() event.ATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return processor.ProcessChannelMessage(input, data)
	}
}

// InteractionHandler 处理内联交互事件
func InteractionHandler() event.InteractionEventHandler {
	return func(event *dto.WSPayload, data *dto.WSInteractionData) error {
		fmt.Println(data)
		return processor.ProcessInlineSearch(data)
	}
}

// GroupATMessageEventHandler 实现处理 at 消息的回调
func GroupATMessageEventHandler() event.GroupATMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGroupATMessageData) error {
		input := strings.ToLower(message.ETLInput(data.Content))
		return processor.ProcessGroupMessage(input, data)
	}
}

// C2CMessageEventHandler 实现处理 at 消息的回调
func C2CMessageEventHandler() event.C2CMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CMessageData) error {
		return processor.ProcessC2CMessage(string(event.RawMessage), data)
	}
}

// C2CFriendEventHandler 实现处理好友关系变更的回调
func C2CFriendEventHandler() event.C2CFriendEventHandler {
	return func(event *dto.WSPayload, data *dto.WSC2CFriendData) error {
		fmt.Println(data)
		return processor.ProcessFriend(string(event.Type), data)
	}
}

// GuildEventHandler 处理频道事件
func GuildEventHandler() event.GuildEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildData) error {
		fmt.Println(data)
		return nil
	}
}

// ChannelEventHandler 处理子频道事件
func ChannelEventHandler() event.ChannelEventHandler {
	return func(event *dto.WSPayload, data *dto.WSChannelData) error {
		fmt.Println(data)
		return nil
	}
}

// GuildMemberEventHandler 处理成员变更事件
func GuildMemberEventHandler() event.GuildMemberEventHandler {
	return func(event *dto.WSPayload, data *dto.WSGuildMemberData) error {
		fmt.Println(data)
		return nil
	}
}

// GuildDirectMessageHandler 处理频道私信事件
func GuildDirectMessageHandler() event.DirectMessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSDirectMessageData) error {
		fmt.Println(data)
		return nil
	}
}

// GuildMessageHandler 处理消息事件
func GuildMessageHandler() event.MessageEventHandler {
	return func(event *dto.WSPayload, data *dto.WSMessageData) error {
		fmt.Println(data)
		return nil
	}
}
