package main

import (
	"fmt"

	"gocv.io/x/gocv"
)

var images = make(chan gocv.Mat)

func FaceDetect() {
	xmlFile := "./haarcascade_frontalcatface.xml"

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	for {
		select {
		case img := <-images:
			if img.Empty() {
				continue
			}
			rects := classifier.DetectMultiScale(img)
			if len(rects) >= 1 {
				fmt.Println("Кол-во Морд: ", len(rects))
				fmt.Println(rects)
			}
		}
	}
}

func main() {
	// open webcam
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	// open display window
	window := gocv.NewWindow("Test Detector")
	defer window.Close()

	// prepare image matrix
	img := gocv.NewMat()
	imgBuf := gocv.NewMat()
	xxxBuf := gocv.NewMat()

	defer img.Close()
	defer imgBuf.Close()
	defer xxxBuf.Close()

	//распознавание морды в кадре
	go FaceDetect()

	qr := gocv.NewQRCodeDetector()
	defer qr.Close()

	counter := 0
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", 0)
			return
		}

		counter += 1
		if counter >= 8 {
			images <- img
			counter = 0
		}

		if img.Empty() {
			continue
		}

		// detect QR
		qrData := qr.DetectAndDecode(img, &imgBuf, &xxxBuf)
		fmt.Printf("Eсть QR: %s - Telemetry imgBuf-%d, xxxBuf-%d\n", qrData, imgBuf.Size(), xxxBuf.Size())

		//blue := color.RGBA{0, 0, 255, 0}
		//for _, r := range img {
		//	gocv.Rectangle(&img, r, blue, 3)
		//
		//	size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
		//	pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
		//	gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		//}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(5) >= 0 {
			break
		}
	}
}
