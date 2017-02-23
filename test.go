// package main

// import (
// 	"fmt"
// )

// type ImgRGB struct {
// 	R int
// 	G int
// 	B int
// }

// // 图像二维数组
// ImgArray := [][]ImgRGB{}

// func main() {
// 	var arr ImgArray
// 	var it ImgRGB
// 	for y := 0; y < 3; y++ {
// 		for x := 0; x < 3; x++ {
// 			r := x + 1
// 			g := x + 2
// 			b := x + 3
// 			// arr[x][y].R = r
// 			// arr[x][y].G = g
// 			// arr[x][y].B = b
// 			it.R = r
// 			it.G = g
// 			it.B = b
// 			// fmt.Println(arr[x][y])
// 			arr = append(arr, it)
// 		}
// 	}
// 	fmt.Println(arr)
// }

// package main

// import (
// 	"fmt"
// )

// type ImgRGB struct {
// 	R int
// 	G int
// 	B int
// }

// type ImgArray [][]ImgRGB

// func main() {
// 	//mapResults := make(map[int]string)
// 	var arrResults ImgArray
// 	for i := 0; i < 5; i++ {
// 		//valueStr := fmt.Sprintf("this is %d", i)
// 		//mapResults[i] = valueStr
// 		var tmpArr []ImgRGB
// 		for j := 0; j < 5; j++ {
// 			var it ImgRGB
// 			it.R = 1
// 			it.G = 0
// 			it.B = 1
// 			tmpArr = append(tmpArr, it)
// 		}
// 		arrResults = append(arrResults, tmpArr)
// 	}
// 	//fmt.Println(mapResults)
// 	fmt.Println(arrResults)
// }

// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/jpeg"
// 	"os"
// )

// func init() {
// 	// damn important or else At(), Bounds() functions will
// 	// caused memory pointer error!!
// 	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
// }

// func main() {
// 	imgfile, err := os.Open("./te.jpeg")

// 	if err != nil {
// 		fmt.Println("img.jpg file not found!")
// 		os.Exit(1)
// 	}

// 	defer imgfile.Close()

// 	// get image height and width with image/jpeg
// 	// change accordinly if file is png or gif

// 	imgCfg, _, err := image.DecodeConfig(imgfile)

// 	if err != nil {
// 		fmt.Println(err)
// 		os.Exit(1)
// 	}

// 	width := imgCfg.Width
// 	height := imgCfg.Height

// 	fmt.Println("Width : ", width)
// 	fmt.Println("Height : ", height)

// 	// we need to reset the io.Reader again for image.Decode() function below to work
// 	// otherwise we will  - panic: runtime error: invalid memory address or nil pointer dereference
// 	// there is no build in rewind for io.Reader, use Seek(0,0)
// 	imgfile.Seek(0, 0)

// 	// get the image
// 	img, _, err := image.Decode(imgfile)

// 	fmt.Println(img.At(10, 10).RGBA())
// 	for y := 0; y < height; y++ {
// 		for x := 0; x < width; x++ {
// 			r, g, b, a := img.At(x, y).RGBA()
// 			fmt.Printf("[X : %d Y : %v] R : %v, G : %v, B : %v, A : %v  \n", x, y, r, g, b, a)
// 		}
// 	}

// }

// package main

// import (
// 	"fmt"
// 	"image"
// 	"image/jpeg"
// 	"os"
// )

// func main() {
// 	// Read image from file that already exists
// 	existingImageFile, err := os.Open("te.jpeg")
// 	if err != nil {
// 		// Handle error
// 	}
// 	defer existingImageFile.Close()

// 	// Calling the generic image.Decode() will tell give us the data
// 	// and type of image it is as a string. We expect "png"
// 	imageData, imageType, err := image.Decode(existingImageFile)
// 	if err != nil {
// 		// Handle error
// 	}
// 	fmt.Println(imageData)
// 	fmt.Println(imageType)

// 	// We only need this because we already read from the file
// 	// We have to reset the file pointer back to beginning
// 	existingImageFile.Seek(0, 0)

// 	// Alternatively, since we know it is a png already
// 	// we can call png.Decode() directly
// 	loadedImage, err := jpeg.Decode(existingImageFile)
// 	if err != nil {
// 		// Handle error
// 	}
// 	fmt.Println(loadedImage)
// }

package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math"
	"os"
)

// kf3 is a generic convolution 3x3 kernel filter that operatates on
// images of type image.Gray from the Go standard image library.
func kf3(k *[9]float64, src, dst *image.Gray) {
	for y := src.Rect.Min.Y; y < src.Rect.Max.Y; y++ {
		for x := src.Rect.Min.X; x < src.Rect.Max.X; x++ {
			var sum float64
			var i int
			for yo := y - 1; yo <= y+1; yo++ {
				for xo := x - 1; xo <= x+1; xo++ {
					if (image.Point{xo, yo}).In(src.Rect) {
						sum += k[i] * float64(src.At(xo, yo).(color.Gray).Y)
					} else {
						sum += k[i] * float64(src.At(x, y).(color.Gray).Y)
					}
					i++
				}
			}
			dst.SetGray(x, y,
				color.Gray{uint8(math.Min(255, math.Max(0, sum)))})
		}
	}
}

// var blur = [9]float64{
// 	1. / 9, 1. / 9, 1. / 9,
// 	1. / 9, 1. / 9, 1. / 9,
// 	1. / 9, 1. / 9, 1. / 9}

var blur = []float64{
	-1, -1, -1,
	-1, 9, -1,
	-1, -1, -1}

// blurY example function applies blur kernel to Y channel
// of YCbCr image using generic kernel filter function kf3
func blurY(src *image.YCbCr) *image.YCbCr {
	dst := *src

	// catch zero-size image here
	if src.Rect.Max.X == src.Rect.Min.X || src.Rect.Max.Y == src.Rect.Min.Y {
		return &dst
	}

	// pass Y channels as gray images
	srcGray := image.Gray{src.Y, src.YStride, src.Rect}
	dstGray := srcGray
	dstGray.Pix = make([]uint8, len(src.Y))
	kf3(&blur, &srcGray, &dstGray) // call generic convolution function

	// complete result
	dst.Y = dstGray.Pix                   // convolution result
	dst.Cb = append([]uint8{}, src.Cb...) // Cb, Cr are just copied
	dst.Cr = append([]uint8{}, src.Cr...)
	return &dst
}

func main() {
	// Example file used here is Lenna100.jpg from the task "Percentage
	// difference between images"
	f, err := os.Open("te.jpeg")
	if err != nil {
		fmt.Println(err)
		return
	}
	img, err := jpeg.Decode(f)
	if err != nil {
		fmt.Println(err)
		return
	}
	f.Close()
	y, ok := img.(*image.YCbCr)
	if !ok {
		fmt.Println("expected color jpeg")
		return
	}
	f, err = os.Create("blur.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = jpeg.Encode(f, blurY(y), &jpeg.Options{90})
	if err != nil {
		fmt.Println(err)
	}
}
