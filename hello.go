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
	runtime.GOMAXPROCS(runtime.NumCPU())  // allocate all the CPU's!
	err := os.MkdirAll(cacheFolder, 0755) // create cache dir if it hasn't been created already
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/makeImg", blurImage) // set up route handler
	http.ListenAndServe(":6060", nil)      // start listening to port 6060 - http://golang.org/pkg/net/http/#ListenAndServe
}

func blurImage(rw http.ResponseWriter, r *http.Request) {

	// Note: Panics are bad for servers.
	// In live settings, find other ways to return error

	// headers are set via rw.Header().Set(key, value)

	imgurl := r.URL.Query().Get("img") // reads the GET-parameter "img"

	cacheName := genCacheName(imgurl)

	file, err := os.OpenFile(cacheName, os.O_RDONLY, 0666) // open cache file for reading
	defer file.Close()                                     // Kill file handler when script ends

	if err == nil {
		// found cache
		_, err = io.Copy(rw, file) // write to output
		if err != nil {
			panic(err)
		}
		return
	}

	resp, err := http.Get(imgurl) // read image file to an http response object
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
	buff := &bytes.Buffer{}                // http://golang.org/pkg/bytes/#Buffer
	err = jpeg.Encode(buff, dstImage, nil) // stuff data in buffer
	if err != nil {
		panic(err)
	}

	cacheFile, err := os.OpenFile(cacheName, os.O_CREATE|os.O_TRUNC, 0755) // open cache file for write, create if not exists, truncate if exists, perms 0755
	defer cacheFile.Close()                                                // Kill file handler when script ends
	if err != nil {
		panic(err)
	}

	multiWriter := io.MultiWriter(cacheFile, rw)          // can write two outputs at once, no sweat - http://golang.org/pkg/io/#MultiWriter
	if _, err := io.Copy(multiWriter, buff); err != nil { // works almost all of the time
		panic(err)
	}
}

func genCacheName(imageUrl string) string {
	b := md5.Sum([]byte(imageUrl))                            // todo: Make sense of this line. Something something hash string tada.
	return path.Join(cacheFolder, fmt.Sprintf("%x.cache", b)) // ditto on the concatenation function here
}
