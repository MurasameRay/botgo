package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// 定义命令注册表
var CommandRegistry = map[string]func(string) string{
	"/hello": helloHandler,
	"/time":  timeHandler,
	"/test":  gptHandler,
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
