package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/chromedp/chromedp"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	var buf = string(Screenshot(r.URL.Query().Get("url")))
	fmt.Fprintf(w, buf)
}

func Screenshot(urlstr string) []byte {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.ExecPath("/headless-shell/headless-shell"),
	}
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// create context
	ctx, cancel := chromedp.NewContext(
		allocCtx,
	)
	defer cancel()

	// capture screenshot of an element
	var buf []byte

	// capture entire browser viewport, returning png with quality=90
	if err := chromedp.Run(ctx, fullScreenshot(urlstr, 90, &buf)); err != nil {
		log.Fatal(err)
	}

	log.Printf("wrote elementScreenshot.png and fullScreenshot.png")

	return buf
}

// fullScreenshot takes a screenshot of the entire browser viewport.
// Note: chromedp.FullScreenshot overrides the device's emulation settings. Reset
func fullScreenshot(urlstr string, quality int, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.EmulateViewport(1200, 630),
		chromedp.FullScreenshot(res, quality),
	}
}
