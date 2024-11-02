package main

import (
	"bufio"
	"fmt"
	"os"
)

func GetUrlFromFile() ([]string, error) {
	// 打开文件
	file, err := os.Open("image_list.txt") // 替换为你的文件路径
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil, err
	}
	defer file.Close() // 确保在函数结束时关闭文件

	var lines []string

	// 使用 bufio.Scanner 逐行读取文件
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// 检查读取过程中是否发生错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}

	// 打印读取的内容
	for _, line := range lines {
		fmt.Println(line)
	}
	return lines, nil
}
