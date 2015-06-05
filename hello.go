/*
	Via https://golang.org/doc/install
*/

package main

import (
	"image/jpeg"
	"net/http"
	"runtime"
)
import "github.com/disintegration/imaging"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/makeImg", blurImage)
	http.ListenAndServe(":6060", nil)
}

func blurImage(rw http.ResponseWriter, r *http.Request) {

	//osr, err := os.OpenFile("file.jpg", os.O_RDONLY, 0666)

	resp, err := http.Get(r.URL.Query().Get("img"))

	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	srcImage, err := jpeg.Decode(resp.Body)

	if err != nil {
		panic(err)
	}

	dstImage := imaging.Blur(srcImage, 28.5)
	//err = imaging.Save(dstImage, "file2.jpg")

	err = jpeg.Encode(rw, dstImage, nil)

	if err != nil {
		panic(err)
	}

}
