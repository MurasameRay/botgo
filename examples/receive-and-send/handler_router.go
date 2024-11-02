package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// 定义命令注册表
var CommandRegistry = map[string]func(string) string{
	"/hello": helloHandler,
	"/time":  timeHandler,
	"/test":  gptHandler,
	"/image": imageHandler,
}

func imageHandler(input string) string {
	index := rand.Int() % 1
	imageList := []string{
		"https://pic2.zhimg.com/v2-738113183a46d8cc5494e269e1356c1d_r.jpg",
		"https://gd-hbimg.huaban.com/a33494c4935115a184a2f74d265e5b4cabd47f20584ea-597UsD_fw1200webp",
		"https://gd-hbimg.huaban.com/c2204b29fd5b1e37b8ea2198c20ab9a8ae017262bf430-QLd98Z",
		"https://gd-hbimg.huaban.com/cc45d37558f6ee802a358d6e76df3539ae7351523c8eb-u9jU9X_fw1200",
	}
	return imageList[index]
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
