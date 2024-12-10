
package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/jpeg"
	"image/png"
	"os"
)

func OverlayImages(overlayPath, backgroundPath, outputPath string) (bool, error) {

	// 指定图片的目录和文件名
	// dir := "picture"
	// backgroundImage := "background.jpg"
	// overlayImage := "overlay.png"
	// outputImage := "output.png"
	// dir, err := os.Getwd()
	// if err != nil {
	// 	fmt.Println("Error:", err)
	// 	return
	// }
	// fmt.Println("Current directory:", dir)

	// 构造图片文件的完整路径
	// backgroundPath := filepath.Join(dir, backgroundImage)
	// overlayPath := filepath.Join(dir, overlayImage)
	// outputPath := filepath.Join(dir, outputImage)

	// 打开底层图片
	baseImgFile, err := os.Open(backgroundPath)
	if err != nil {
		return false, fmt.Errorf("could not open background image: %v", err)
	}
	defer baseImgFile.Close()

	// 解码底层图片
	baseImg, _, err := image.Decode(baseImgFile)
	if err != nil {
		return false, fmt.Errorf("could not decode background image: %v", err)
	}

	// 打开覆盖图片
	overlayImgFile, err := os.Open(overlayPath)
	if err != nil {
		return false, fmt.Errorf("could not open overlay image: %v", err)
	}
	defer overlayImgFile.Close()

	// 解码覆盖图片
	overlayImg, err := png.Decode(overlayImgFile)
	if err != nil {
		return false, fmt.Errorf("could not decode foreground image: %v", err)
	}

	// 设置缩放后的宽度和高度（例如，将宽度缩小到200像素，高度按比例缩放）
	// newWidth := 750
	// originalAspectRatio := float64(overlayImg.Bounds().Dy()) / float64(overlayImg.Bounds().Dx())
	// newHeight := int(float64(newWidth) * originalAspectRatio)

	// 等比缩放图片
	// m := resize.Resize(uint(newWidth), uint(newHeight), overlayImg, resize.Lanczos3)

	// 创建一个新的RGBA图像，用于存储最终合成的图像
	outputImg := image.NewRGBA(baseImg.Bounds())

	// 将底层图片绘制到outputImg上
	draw.Draw(outputImg, outputImg.Bounds(), baseImg, image.Point{0, 0}, draw.Src)

	// 将覆盖图片绘制到outputImg上
	overlayPosition := image.Pt(25, 50) // 覆盖图片在底层图片上的位置（例如50, 50位置）
	draw.Draw(outputImg, overlayImg.Bounds().Add(overlayPosition), overlayImg, image.Point{0, 0}, draw.Over)

	// 创建输出文件
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return false, fmt.Errorf("could not create output file: %v", err)
	}
	defer outputFile.Close()

	// 将最终的图像写入输出文件
	err = png.Encode(outputFile, outputImg)
	if err != nil {
		return false, fmt.Errorf("could not encode output image: %v", err)
	}

	fmt.Println("Image overlay completed successfully and saved as output.png")
	return true, nil
}

