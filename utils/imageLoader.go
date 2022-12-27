// For loading images that have allready been preprocessed
package main

import (
    // "bytes"
	"embed"
    "image"
    "image/draw"
    "image/gif"
    "image/color/palette"
    "log"
    // "fmt"
    // "os"
    "time"

    "periph.io/x/conn/i2c/i2creg"
    "periph.io/x/devices/ssd1306"
    // "periph.io/x/devices/ssd1306/image1bit"
    // "periph.io/x/conn/physic"
	"periph.io/x/host"
    

    "github.com/nfnt/resize"
)


//go:embed images
var image_path embed.FS



func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}
 
    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")

    // Open a handle to a ssd1306 connected on the I²C bus:
    dev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
    
    // imgBunnyNRGBA, err := gif.Decode(bytes.NewReader(bunny))
    // if err != nil {
    //     log.Fatal(err)
    // }
    // img := convertAndResizeAndCenter(128, 64, imgBunnyNRGBA)

    path := "images/Avatar.gif"
    f, err := image_path.Open(path)
	if err != nil {
		log.Fatal(err)
	}
    
    // img, err := gif.Decode(bytes.NewReader(f))
    // buf := new(bytes.Buffer)
    // err = gif.Encode(buf, img, nil)
    // send_s3 := buf.Bytes()
    // err = os.WriteFile("/tmp/dat1", send_s3, 0644)
    if err != nil {
        log.Fatal(err)
    }
    // center := imgBunny1bit.Bounds()
    
    // draw.Src.Draw(imgBunny1bitLarge, center.Add(image.Point{X: (128 - center.Dx()) / 2}), imgBunny1bit, image.Point{})
	// dev.Draw(img.Bounds(), img, image.Point{})
    // imgClear := make([]byte, 128*64/8)

   
    //  _, err = dev.Write(imgClear);
	
    if err != nil {
        log.Fatal(err)
    }
    g, err := gif.DecodeAll(f)
    log.Println(g.Delay)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }

    
    // Converts every frame to image.Gray and resize them:
    // outGif := &gif.GIF{}
    // imgs := make([]*image.Paletted, len(g.Image))
    // for i := range g.Image {
    //     imgs[i] = convertAndResizeAndCenter(128, 64, g.Image[i])
    //     outGif.Image = append(outGif.Image, imgs[i])
    //     outGif.Delay = append(outGif.Delay, 2)
    // }
    
    // log.Println(outGif.Delay)

    // buf := new(bytes.Buffer)
    // err = gif.EncodeAll(buf, outGif)
    // send_s3 := buf.Bytes()
    // err = os.WriteFile("/home/pi/go/src/pi/helloworld/images/wave.gif", send_s3, 0644)


    // Display the frames in a loop:
    for i := 1; ; i++ {
        index := i % len(g.Image)
        c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
        img := g.Image[index]
        dev.Draw(img.Bounds(), img, image.Point{})
        // if i == len(g.Image) + 1 {
        //     break
        // }
        <-c
    }
}