package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

func GetImageURLs(folderPath string, urlPrefix string) ([]string, error) {
	var imageURLs []string

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 检查是否是文件且后缀为 .jpg 或 .png
		if !info.IsDir() {
			if strings.HasSuffix(strings.ToLower(info.Name()), ".jpg") || strings.HasSuffix(strings.ToLower(info.Name()), ".png") {
				// 计算相对路径
				relPath, err := filepath.Rel(folderPath, path)
				if err != nil {
					return err
				}

				// 拼接 URL
				imageURL := fmt.Sprintf("%s/%s", urlPrefix, relPath)
				imageURLs = append(imageURLs, imageURL)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return imageURLs, nil
}
