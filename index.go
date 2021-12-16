package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func handler(w http.ResponseWriter, r *http.Request) {
	println(r.URL.Query().Get("url"))
	var fileBytes = Screenshot(r.URL.Query().Get("url"))

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(fileBytes)
	return
}

func Screenshot(urlstr string) []byte {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(urlstr, 100, &buf)); err != nil {
		log.Fatal(err)
	}

	return buf
}

// fullScreenshot takes a screenshot of the entire browser viewport.
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		//chromedp.EmulateScale(2)
		chromedp.EmulateViewport(1200, 630),
		chromedp.WaitVisible(`#logo`, chromedp.ByID),
		chromedp.FullScreenshot(res, quality),
	}
}

func main() {
	port := ":3000"

	http.HandleFunc("/api/screenshot", handler)
	fmt.Println("Listening in port " + port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
