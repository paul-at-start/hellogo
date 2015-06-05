/*
	Via https://golang.org/doc/install
*/

package main

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
)
import "github.com/disintegration/imaging"

const cacheFolder = "imgcache"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	err := os.MkdirAll(cacheFolder, 0755)
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/makeImg", blurImage)
	http.ListenAndServe(":6060", nil)
}

func blurImage(rw http.ResponseWriter, r *http.Request) {

	// Note: Panics are bad.
	// In live settings, find other ways to return error

	// headers are set via rw.Header().Set(key, value)

	imgurl := r.URL.Query().Get("img")

	cacheName := genCacheName(imgurl)

	// open cache for reading
	file, err := os.OpenFile(cacheName, os.O_RDONLY, 0666)
	defer file.Close()

	if err == nil {
		// found cache
		_, err = io.Copy(rw, file) // write to output
		if err != nil {
			panic(err)
		}
		return
	}

	resp, err := http.Get(imgurl)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(err)
	}

	srcImage, err := jpeg.Decode(resp.Body)

	if err != nil {
		panic(err)
	}

	dstImage := imaging.Blur(srcImage, 28.5)

	// write to cache

	buff := &bytes.Buffer{}
	err = jpeg.Encode(buff, dstImage, nil)
	if err != nil {
		panic(err)
	}

	cacheFile, err := os.OpenFile(cacheName, os.O_CREATE|os.O_TRUNC, 0755)
	defer cacheFile.Close()
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(cacheFile, rw)
	if _, err := io.Copy(multiWriter, buff); err != nil {
		panic(err)
	}
}

func genCacheName(imageUrl string) string {
	b := md5.Sum([]byte(imageUrl))
	return path.Join(cacheFolder, fmt.Sprintf("%x.cache", b))
}
