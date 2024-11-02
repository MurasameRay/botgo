package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// 定义命令注册表
var CommandRegistry = map[string]func(string) string{
	"/hello": helloHandler,
	"/time":  timeHandler,
	"/image": imageHandler,
	"":       gptHandler,
}

func imageHandler(input string) string {

	imageList := []string{
		"https://pic2.zhimg.com/v2-738113183a46d8cc5494e269e1356c1d_r.jpg",
		"https://pic2.zhimg.com/v2-53efed13d3b03ea0a831e94035b32f9b_1440w.jpg",
		"https://i0.hdslb.com/bfs/new_dyn/621d179e761b9ddc062f2d5e57516b37166383626.jpg",
		"https://konachan.com/sample/7e64428813fb7119530f388777f7d9ba/Konachan.com%20-%20382933%20sample.jpg",
		//"https://gd-hbimg.huaban.com/a33494c4935115a184a2f74d265e5b4cabd47f20584ea-597UsD_fw1200webp",
		//"https://gd-hbimg.huaban.com/c2204b29fd5b1e37b8ea2198c20ab9a8ae017262bf430-QLd98Z",
		//"https://gd-hbimg.huaban.com/cc45d37558f6ee802a358d6e76df3539ae7351523c8eb-u9jU9X_fw1200",
	}
	index := rand.Int() % len(imageList)
	//bytes, err := DownloadImage(imageList[index])
	//if err != nil {
	//	return err.Error()
	//}
	//return string(bytes)
	return imageList[index]
}

// DownloadImage 下载指定 URL 的图片并返回二进制数据
func DownloadImage(url string) ([]byte, error) {
	// 发送 GET 请求
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to download image: %w", err)
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download image: status code %d", resp.StatusCode)
	}

	// 读取响应体
	imageData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	return imageData, nil
}

// 处理命令的函数
func ProcessCommand(command string) string {
	for prefix, handler := range CommandRegistry {
		if strings.HasPrefix(command, prefix) {
			// 去掉前缀并提取消息
			message := strings.TrimSpace(strings.TrimPrefix(command, prefix))
			return handler(message)
		}
	}
	return "Unknown command"
}

// 定义处理函数
func helloHandler(message string) string {
	return fmt.Sprintf("hello, %s", message)
}

func timeHandler(message string) string {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("The current time is: %s", currentTime)
}

func gptHandler(message string) string {
	url := "http://127.0.0.1:23333/ask_chatgpt?word=" + message
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return err.Error()
	}
	return string(body)
}
