package main

import (
	"fmt"
	"github.com/bregydoc/gtranslate"
	"image"
	"image/color"
	"io/ioutil"
	"unicode/utf8"

	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

const (
	FontSize     = 20
	SquareHeight = 200
)

func main() {
	HandleUserImage()
}

// HandleUserImage paste user image onto background
func HandleUserImage() (string, error) {
	m, err := imaging.Open("target.png")
	if err != nil {
		fmt.Println("open file failed", err)
		return "", err
	}

	bm, err := imaging.Open("bg.jpg")
	if err != nil {
		fmt.Println("open file failed", err)
		return "", err
	}

	// 图片按比例缩放
	dst := imaging.Resize(m, bm.Bounds().Max.X, SquareHeight, imaging.Lanczos)
	// 将图片粘贴到背景图的固定位置
	result := imaging.Overlay(bm, dst, image.Pt(0, bm.Bounds().Max.Y-SquareHeight), 0.3)
	writeOnImage(result)

	fileName := fmt.Sprintf("%d.jpg", 1234)
	err = imaging.Save(result, fileName)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func writeOnImage(target *image.NRGBA) {
	c := freetype.NewContext()

	c.SetDPI(256)
	c.SetClip(target.Bounds())
	c.SetDst(target)
	c.SetHinting(font.HintingFull)

	// 设置文字颜色、字体、字大小
	c.SetSrc(image.NewUniform(color.RGBA{R: 220, G: 220, B: 220, A: 220}))
	c.SetFontSize(FontSize)
	fontFam, err := getFontFamily()
	if err != nil {
		fmt.Println("get font family error")
	}
	c.SetFont(fontFam)

	drawStr := "你好，世界"

	// 获取字体的尺寸大小
	fixed := c.PointToFixed(FontSize)
	// fixed.Ceil() 字体大小
	// utf8.RuneCountInString(drawStr) 获取字符串的实际大小，而不是以byte算
	pt := freetype.Pt(target.Rect.Max.X/2-(utf8.RuneCountInString(drawStr)/2)*fixed.Ceil(), target.Rect.Max.Y-SquareHeight+SquareHeight/2)

	fontPin, err := c.DrawString(drawStr, pt)
	if err != nil {
		fmt.Printf("draw error: %v \n", err)
		return
	}

	// 翻译文字
	translated, err := gtranslate.TranslateWithFromTo(
		drawStr,
		gtranslate.FromTo{
			From: "zh-cn",
			To:   "ja",
		},
	)
	if err != nil {
		fmt.Printf("translate error: %v \n", err)
		return
	}

	_, err = c.DrawString(translated, freetype.Pt(target.Rect.Max.X/2-(utf8.RuneCountInString(translated)/2)*fixed.Ceil(), fontPin.Y.Ceil()+fixed.Ceil()))
	if err != nil {
		fmt.Printf("draw error: %v \n", err)
		return
	}
}

// 获取字符集，仅调用一次
func getFontFamily() (*truetype.Font, error) {
	// 这里需要读取中文字体，否则中文文字会变成方格
	fontBytes, err := ioutil.ReadFile("msyh.ttc")
	if err != nil {
		fmt.Println("read file error:", err)
		return &truetype.Font{}, err
	}

	f, err := freetype.ParseFont(fontBytes)
	if err != nil {
		fmt.Println("parse font error:", err)
		return &truetype.Font{}, err
	}

	return f, err
}
