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
//var CommandRegistry = map[string]func(string) string{
//	"/hello": helloHandler,
//	"/time":  timeHandler,
//	"/image": imageHandler,
//	"/gpt":   gptHandler,
//	"":       gptHandler,
//}

// 定义命令注册 不随机匹配，按照从上到下的顺序
var CommandRegistry = []struct {
	Prefix  string
	Handler func(string) string
}{
	{"/hello", helloHandler},
	{"/time", timeHandler},
	{"/image", imageHandler},
	//{"/gpt", gptHandler},
	//{"", gptHandler},
	//{"", helloHandler},
}

func imageHandler(input string) string {
	//folderPath := "/www/wwwroot/blog.xiaocongyu.com/wp-content/uploads" // 替换为您的文件夹路径
	folderPath := "/root/images" // 替换为您的文件夹路径
	//urlPrefix := "http://blog.xiaocongyu.com/wp-content/uploads"
	urlPrefix := "http://file.xiaocongyu.com"
	if strings.Contains(input, "ba") || strings.Contains(input, "blue archive") ||
		strings.Contains(input, "blue") || strings.Contains(input, "档案") {
		folderPath += "/blue_archive"
		urlPrefix += "/blue_archive"
	}

	imageList, _ := GetImageURLs(folderPath, urlPrefix)
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
	for _, handler := range CommandRegistry {
		if strings.HasPrefix(command, handler.Prefix) {
			// 去掉前缀并提取消息
			fmt.Println("Original command:", command)
			message := strings.TrimSpace(strings.TrimPrefix(command, handler.Prefix))
			return handler.Handler(message)
		}
	}
	return "Unknown command"
}

// 定义处理函数
func helloHandler(message string) string {
	if message == "" {
		return "老师你好，我是阿罗娜～有什么可以帮你的吗"
	}
	return message
}

func timeHandler(message string) string {
	// 加载北京时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return fmt.Sprintf("Error loading location: %v", err)
	}

	// 获取当前时间并设置为北京时区
	currentTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
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
