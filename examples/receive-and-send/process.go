package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tencent-connect/botgo/dto"
	"github.com/tencent-connect/botgo/openapi"
)

// Processor is a struct to process message
type Processor struct {
	api openapi.OpenAPI
}

// ProcessChannelMessage is a function to process message
func (p Processor) ProcessChannelMessage(input string, data *dto.WSATMessageData) error {
	msg := generateDemoMessage(input, dto.Message(*data))
	if err := p.sendChannelReply(context.Background(), data.ChannelID, msg); err != nil {
		_ = p.sendChannelReply(context.Background(), data.GroupID, genErrMessage(dto.Message(*data), err))
	}
	return nil
}

// ProcessInlineSearch is a function to process inline search
func (p Processor) ProcessInlineSearch(interaction *dto.WSInteractionData) error {
	if interaction.Data.Type != dto.InteractionDataTypeChatSearch {
		return fmt.Errorf("interaction data type not chat search")
	}
	search := &dto.SearchInputResolved{}
	if err := json.Unmarshal(interaction.Data.Resolved, search); err != nil {
		log.Println(err)
		return err
	}
	if search.Keyword != "test" {
		return fmt.Errorf("resolved search key not allowed")
	}
	searchRsp := &dto.SearchRsp{
		Layouts: []dto.SearchLayout{
			{
				LayoutType: 0,
				ActionType: 0,
				Title:      "内联搜索",
				Records: []dto.SearchRecord{
					{
						Cover: "https://pub.idqqimg.com/pc/misc/files/20211208/311cfc87ce394c62b7c9f0508658cf25.png",
						Title: "内联搜索标题",
						Tips:  "内联搜索 tips",
						URL:   "https://www.qq.com",
					},
				},
			},
		},
	}
	body, _ := json.Marshal(searchRsp)
	if err := p.api.PutInteraction(context.Background(), interaction.ID, string(body)); err != nil {
		log.Println("api call putInteractionInlineSearch  error: ", err)
		return err
	}
	return nil
}

func genErrMessage(data dto.Message, err error) *dto.MessageToCreate {
	return &dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   fmt.Sprintf("处理异常:%v", err),
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
		MsgID: data.ID,
	}
}

// ProcessGroupMessage 回复群消息
func (p Processor) ProcessGroupMessage(input string, data *dto.WSGroupATMessageData) error {
	msg := generateDemoMessage(input, dto.Message(*data))
	if err := p.sendGroupReply(context.Background(), data.GroupID, msg); err != nil {
		_ = p.sendGroupReply(context.Background(), data.GroupID, genErrMessage(dto.Message(*data), err))
	}

	return nil
}

// ProcessC2CMessage 回复C2C消息
func (p Processor) ProcessC2CMessage(input string, data *dto.WSC2CMessageData) error {
	userID := ""
	if data.Author != nil && data.Author.ID != "" {
		userID = data.Author.ID
	}
	msg := generateDemoMessage(input, dto.Message(*data))
	if err := p.sendC2CReply(context.Background(), userID, msg); err != nil {
		_ = p.sendC2CReply(context.Background(), userID, genErrMessage(dto.Message(*data), err))
	}
	return nil
}

func generateDemoMessage(input string, data dto.Message) *dto.MessageToCreate {
	log.Printf("收到指令: %+v", input)
	msg := ProcessCommand(input)

	//msg := ""
	//if len(input) > 0 {
	//	msg += "收到:" + input
	//}
	//for _, _v := range data.Attachments {
	//	msg += ",收到文件类型:" + _v.ContentType
	//}
	msgType := dto.TextMsg
	response := &dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   msg,
		MsgType:   msgType,
		MessageReference: &dto.MessageReference{
			// 引用这条消息
			MessageID:             data.ID,
			IgnoreGetMessageError: true,
		},
		MsgID: data.ID,
	}
	if strings.HasPrefix(msg, "http") {
		file, err := UploadFile(data.GroupID, 1, msg, false)

		if err != nil {
			response.Content = err.Error()
			return response
		}
		response.MsgType = dto.RichMediaMsg
		response.Media = &dto.MediaInfo{
			FileInfo: []byte(file.FileInfo),
		}
		response.Content = "图片效果"
	}
	return response
}

// ProcessFriend 处理 c2c 好友事件
func (p Processor) ProcessFriend(wsEventType string, data *dto.WSC2CFriendData) error {
	// 请注意，这里是主动推送添加好友事件，后续改为 event id 被动消息
	replyMsg := dto.MessageToCreate{
		Timestamp: time.Now().UnixMilli(),
		Content:   "",
	}
	var content string
	switch strings.ToLower(wsEventType) {
	case strings.ToLower(string(dto.EventC2CFriendAdd)):
		log.Println("添加好友")
		content = fmt.Sprintf("ID为 %s 的用户添加机器人为好友", data.OpenID)
	case strings.ToLower(string(dto.EventC2CFriendDel)):
		log.Println("删除好友")
	default:
		log.Println(wsEventType)
		return nil
	}
	replyMsg.Content = content
	_, err := p.api.PostC2CMessage(
		context.Background(),
		data.OpenID,
		replyMsg,
	)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// FileUploadRequest 定义请求参数结构
type FileUploadRequest struct {
	FileType   int    `json:"file_type"`
	URL        string `json:"url"`
	SrvSendMsg bool   `json:"srv_send_msg"`
	FileData   string `json:"file_data,omitempty"` // 可选字段
}

// FileUploadResponse 定义响应参数结构
type FileUploadResponse struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
	ID       string `json:"id,omitempty"` // 可选字段
}

// UploadFile 上传文件到群聊的函数
func UploadFile(groupOpenID string, fileType int, url string, srvSendMsg bool) (*FileUploadResponse, error) {
	// 请求参数
	// 构建请求 URL
	reqURL := fmt.Sprintf("https://api.sgroup.qq.com/v2/groups/%s/files", groupOpenID)
	method := "POST"

	payload := strings.NewReader(`{
  "file_type": ` + strconv.Itoa(fileType) + `,
  "url": "` + url + `",
  "srv_send_msg": ` + strconv.FormatBool(srvSendMsg) + `
}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, reqURL, payload)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	authToken, _ := tokenSource.Token()
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("appid", "102457514")
	req.Header.Add("Authorization", "QQBot "+authToken.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//fmt.Println(string(body))
	// 解析响应
	var response FileUploadResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	// 检查响应状态
	if response.FileInfo == "" {
		return nil, fmt.Errorf("request failed with body: %s", body)
	}
	return &response, nil
}

//
//// GetAccessToken 从配置文件中获取 Token
//func GetAccessToken(ctx context.Context) (string, error) {
//	// 读取配置文件
//	content, err := os.ReadFile("config.yaml")
//	if err != nil {
//		return "", fmt.Errorf("load config file failed, err: %w", err)
//	}
//
//	// 解析配置文件
//	credentials := &token.QQBotCredentials{}
//	if err = yaml.Unmarshal(content, credentials); err != nil {
//		return "", fmt.Errorf("parse config failed, err: %w", err)
//	}
//
//	// 创建 Token 源
//	tokenSource := token.NewQQBotTokenSource(credentials)
//
//	// 刷新 Access Token
//	if err = token.StartRefreshAccessToken(ctx, tokenSource); err != nil {
//		return "", fmt.Errorf("failed to refresh access token: %w", err)
//	}
//
//	// 获取 Token
//	token, _ := tokenSource.Token() // 假设 Token() 方法返回当前 Token
//	return token.AccessToken, nil
//}
