package main

import (
	"fmt"
	//"image"
	//"image/draw"
	"image/jpeg"
	"os"
)

// 像素单元
type ImgRGB struct {
	R uint32
	G uint32
	B uint32
}

// 图像二维数组
type ImgArray [][]ImgRGB

// 卷积核
//var Filter [][]ImgRGB
// var Fone = [3][3]int{{1, 0, 1}, {0, 1, 0}, {1, 0, 1}}

// var Filter interface {
// 	Fone
// }

func ImgLen() (width int, height int, arr ImgArray) {
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
	// gray := image.NewGray(size)

	// draw.Draw(gray, size, img, img.Bounds().Min, draw.Src)

	width = size.Dx()
	height = size.Dy()

	for y := size.Min.Y; y < size.Max.Y; y++ {
		var tmp []ImgRGB
		for x := size.Min.X; x < size.Max.X; x++ {
			var item ImgRGB
			r, g, b, _ := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			//  Shifting by 12 reduces this tothe range [0, 15].
			item.R = r
			item.G = g
			item.B = b
			tmp = append(tmp, item)
			//fmt.Printf("a[%d][%d] = %d\n", x,y, arr[x][y])
		}
		arr = append(arr, tmp)
	}

	return
}

func main() {
	_, _, re := ImgLen()
	fmt.Println(re)
}
