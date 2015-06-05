/*
	Via https://golang.org/doc/install
*/

package main

import (
	"image/jpeg"
	"net/http"
	"os"
)
import "github.com/disintegration/imaging"

func main() {
	http.HandleFunc("/makeImg", blurImage)
	http.ListenAndServe(":6060", nil)
}

func blurImage(rw http.ResponseWriter, r *http.Request) {

	osr, err := os.OpenFile("file.jpg", os.O_RDONLY, 0666)

	if err != nil {
		panic(err)
	}

	srcImage, err := jpeg.Decode(osr)

	if err != nil {
		panic(err)
	}

	dstImage := imaging.Blur(srcImage, 13.5)
	//err = imaging.Save(dstImage, "file2.jpg")

	err = jpeg.Encode(rw, dstImage, nil)

	if err != nil {
		panic(err)
	}

}
