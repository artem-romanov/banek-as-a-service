package memer

import (
	_ "embed"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

type ImageType string

const (
	Png  ImageType = "png"
	Jpeg ImageType = "jpeg"
)

var strToImageType = map[string]ImageType{
	"png":  Png,
	"jpeg": Jpeg,
}

type Memer struct {
	BottomMargin int
	SideMargin   int
	font         *opentype.Font
}

//go:embed assets/fonts/impact.ttf
var impactTTF []byte

func NewMemer(bottomMargin, sideMargin int) (*Memer, error) {
	f, err := opentype.Parse(impactTTF)
	if err != nil {
		return nil, fmt.Errorf("cant init memer, parse font err: %w", err)
	}

	return &Memer{
		BottomMargin: bottomMargin,
		SideMargin:   sideMargin,
		font:         f,
	}, nil
}

// Warning! Mutates srcImg to avoid double mem allocations.
func (m *Memer) Generate(srcImg image.Image, text string) (image.Image, error) {
	var destImg draw.Image

	switch img := srcImg.(type) {
	case *image.RGBA:
		destImg = img
	case *image.NRGBA:
		// Конвертируем в RGBA (одна копия)
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		destImg = rgba
	default:
		// Для YCbCr и прочего — одна копия в RGBA
		rgba := image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)
		destImg = rgba
	}

	imgW := float64(srcImg.Bounds().Dx())
	imgH := float64(srcImg.Bounds().Dy())

	// 1. Базовый размер: 15% высоты изображения
	baseSize := imgH * 0.15

	face, err := opentype.NewFace(m.font, &opentype.FaceOptions{
		Size:    baseSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("generate err: %w", err)
	}
	defer face.Close()

	maxWidth := imgW - float64(2*m.SideMargin)
	width := font.MeasureString(face, text).Ceil()

	// 2. Если текст не влезает по ширине — уменьшаем
	if float64(width) > maxWidth {
		scale := maxWidth / float64(width)
		newSize := baseSize * scale

		face, _ = opentype.NewFace(m.font, &opentype.FaceOptions{
			Size:    newSize,
			DPI:     72,
			Hinting: font.HintingFull,
		})
	}

	d := font.Drawer{
		Dst:  destImg,
		Src:  image.NewUniform(color.White),
		Face: face,
	}

	// центрируем
	textW := font.MeasureString(face, text).Ceil()
	x := (srcImg.Bounds().Dx() - textW) / 2
	y := srcImg.Bounds().Dy() - m.BottomMargin

	d.Dot = fixed.P(x, y)
	drawTextWithOutline(&d, x, y, text, color.Black, color.White, 2)

	return destImg, nil
}

func FileToImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("open file err", err)
	}
	defer file.Close()

	_, formatStr, err := image.DecodeConfig(file)
	if err != nil {
		return nil, fmt.Errorf("decoding config err: %w", err)
	}

	if _, err := file.Seek(0, 0); err != nil {
		return nil, fmt.Errorf("file to image err: %w", err)
	}

	imageType, ok := strToImageType[formatStr]
	if !ok {
		return nil, errors.New("unknown file format")
	}

	var img image.Image
	switch imageType {
	case Jpeg:
		img, err = jpeg.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("decoding image err: %w", err)
		}
	case Png:
		img, err = png.Decode(file)
		if err != nil {
			return nil, fmt.Errorf("decoding image err: %w", err)
		}
	}

	return img, nil
}

func WriteImageToJpegFile(destImg image.Image, filename string) error {
	destFile, err := os.Create(filename + ".jpg")
	if err != nil {
		log.Fatalln("coudn't create dest file", err)
	}
	defer destFile.Close()

	if err := jpeg.Encode(destFile, destImg, &jpeg.Options{Quality: 95}); err != nil {
		return fmt.Errorf("coudn't write jpeg: %w", err)
	}
	return nil
}

func drawTextWithOutline(d *font.Drawer, x, y int, text string, outlineColor, fillColor color.Color, thickness int) {
	offsets := []image.Point{}
	for dx := -thickness; dx <= thickness; dx++ {
		for dy := -thickness; dy <= thickness; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}
			offsets = append(offsets, image.Point{X: dx, Y: dy})
		}
	}

	// 1. Рисуем обводку
	d.Src = image.NewUniform(outlineColor)
	for _, o := range offsets {
		d.Dot = fixed.P(x+o.X, y+o.Y)
		d.DrawString(text)
	}

	// 2. Рисуем основной текст
	d.Src = image.NewUniform(fillColor)
	d.Dot = fixed.P(x, y)
	d.DrawString(text)
}
