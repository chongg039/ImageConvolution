package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

func Thresholding() {
	f, err := os.Open("./te.jpeg")
	if err != nil {
		return
	}
	defer f.Close()

	img, err := jpeg.Decode(f)
	if err != nil {
		return
	}
	size := img.Bounds()
	gray := image.NewGray(size)

	draw.Draw(gray, size, img, img.Bounds().Min, draw.Src)

	width := size.Dx()
	height := size.Dy()

	// zft := make([]int, 256)
	var zft []int
	var idx int

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			idx = i*height + j
			zft[int(gray.Pix[idx])]++
		}
	}

	fz := GetOSTUThreshold(zft)
	for i := 0; i < len(gray.Pix); i++ {
		if gray.Pix[i] > fz {
			gray.Pix[i] = 255
		} else {
			gray.Pix[i] = 0
		}
	}
	w, _ := os.Create("./gray.jpeg")
	jpeg.Encode(w, gray, nil)
}

func GetOSTUThreshold(HistGram []int) int {
	var Y, Amount int
	var PixelBack, PixelFore, PixelIntegralBack, PixelIntegralFore, PixelIntegral int
	var OmegaBack, OmegaFore, MicroBack, MicroFore, SigmaB, Sigma float64 // 类间方差;
	var MinValue, MaxValue int
	var Threshold int = 0

	for MinValue = 0; MinValue < 256 && HistGram[MinValue] == 0; MinValue++ {

	}
	for MaxValue = 255; MaxValue > MinValue && HistGram[MinValue] == 0; MaxValue-- {

	}
	if MaxValue == MinValue {
		return MaxValue // 图像中只有一个颜色
	}
	if MinValue+1 == MaxValue {
		return MinValue // 图像中只有二个颜色
	}

	for Y = MinValue; Y <= MaxValue; Y++ {
		Amount += HistGram[Y] //  像素总数
	}
	PixelIntegral = 0
	for Y = MinValue; Y <= MaxValue; Y++ {
		PixelIntegral += HistGram[Y] * Y
	}
	SigmaB = -1
	for Y = MinValue; Y < MaxValue; Y++ {
		PixelBack = PixelBack + HistGram[Y]
		PixelFore = Amount - PixelBack
		OmegaBack = float64(PixelBack) / float64(Amount)
		OmegaFore = float64(PixelFore) / float64(Amount)
		PixelIntegralBack += HistGram[Y] * Y
		PixelIntegralFore = PixelIntegral - PixelIntegralBack
		MicroBack = float64(PixelIntegralBack) / float64(PixelBack)
		MicroFore = float64(PixelIntegralFore) / float64(PixelFore)
		Sigma = OmegaBack * OmegaFore * (MicroBack - MicroFore) * (MicroBack - MicroFore)
		if Sigma > SigmaB {
			SigmaB = Sigma
			Threshold = Y
		}
	}
	return Threshold
}

func main() {
	Thresholding()
}
