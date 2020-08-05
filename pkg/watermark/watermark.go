package watermark

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"errors"
	"github.com/disintegration/imaging"
	"github.com/golang/freetype"
	"github.com/pefish/go-watermarker/pkg/internal/asset"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
)

var ErrUnsupportedWatermarkType = errors.New("不支持的类型")

func (w *Watermark) MarkGif(srcFile *os.File, dstPath string) error {
	return ErrUnsupportedWatermarkType
}

func (w *Watermark) markJpgAndPng(srcImg image.Image) (image.Image, error) {
	// 新建背景图层，将源图片全部画上去
	img := image.NewNRGBA(srcImg.Bounds())
	for y := 0; y < img.Bounds().Dy(); y++ {
		for x := 0; x < img.Bounds().Dx(); x++ {
			img.Set(x, y, srcImg.At(x, y))
		}
	}
	//draw.Draw(img, jpgImg.Bounds(), image.White, image.ZP, draw.Over)

	// 新建字体图层
	fontImg := image.NewNRGBA(srcImg.Bounds())

	// 解压出font资源
	fontGzipBytes, err := hex.DecodeString(asset.FontStr)
	if err != nil {
		return nil, err
	}
	var dstBuf bytes.Buffer
	zr, err := gzip.NewReader(bytes.NewReader(fontGzipBytes))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	_, err = io.Copy(&dstBuf, zr)
	if err != nil {
		return nil, err
	}

	fontBytes := dstBuf.Bytes()
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return nil, err
	}
	f := freetype.NewContext()
	f.SetDPI(72)
	f.SetFont(font)
	f.SetFontSize(40)
	f.SetClip(fontImg.Bounds())
	f.SetDst(fontImg)  // 设置写字的目标图层
	f.SetSrc(image.NewUniform(w.textColor))
	pt := freetype.Pt(fontImg.Bounds().Min.X + fontImg.Bounds().Dx()/4, fontImg.Bounds().Min.Y + 60)  // 左上角原点，x轴向下。参数是开始写字的位置
	_, err = f.DrawString(w.text, pt)  // 将字写上去
	if err != nil {
		return nil, err
	}

	// 字体图层旋转30度
	fontImg = imaging.Rotate(fontImg, 20, color.Transparent)  // 右下角为支点旋转

	// 字体图层合并到背景图层
	draw.Draw(img, img.Bounds(), fontImg, image.ZP, draw.Over)


	return img, nil
}

func (w *Watermark) MarkJpg(srcFile *os.File, dstPath string) error {
	jpgImg, err := jpeg.Decode(srcFile)
	if err != nil {
		return err
	}

	img, err := w.markJpgAndPng(jpgImg)
	if err != nil {
		return err
	}
	//保存到新文件中
	newFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	err = jpeg.Encode(newFile, img, &jpeg.Options{Quality: 100})
	if err != nil {
		return err
	}
	return nil
}

func (w *Watermark) MarkPng(srcFile *os.File, dstPath string) error {
	pngImg, err := png.Decode(srcFile)
	if err != nil {
		return err
	}
	img, err := w.markJpgAndPng(pngImg)
	if err != nil {
		return err
	}
	newFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer newFile.Close()
	err = png.Encode(newFile, img)
	if err != nil {
		return err
	}
	return nil
}

type Watermark struct {
	text      string
	textColor color.RGBA
	allowExts map[string]func(srcPath *os.File, dstPath string) error
}

func NewWatermark(text string) *Watermark {
	w := &Watermark{
		text:      text,
		textColor: color.RGBA{  // 透明度为51/255的纯红色。对应是[255,0,0,0.2]
			R: 51,  // 这是预乘值（真正的R*透明度）
			G: 0,
			B: 0,
			A: 51,  // 51/255 = 0.2
		},
	}
	w.allowExts = map[string]func(srcPath *os.File, dstPath string) error{
		".gif": w.MarkGif, ".jpg": w.MarkJpg, ".jpeg": w.MarkJpg, ".png": w.MarkPng,
	}
	return w
}

func (w *Watermark) SetTextColor(textColor color.RGBA) {
	w.textColor = textColor
}

// MarkFile 给指定的文件打上水印
func (w *Watermark) MarkFile(srcPath, dstPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	ext := path.Ext(srcPath)
	markFun, ok := w.allowExts[ext]
	if !ok {
		return ErrUnsupportedWatermarkType
	}
	err = markFun(srcFile, dstPath)
	if err != nil {
		return err
	}
	return nil
}
