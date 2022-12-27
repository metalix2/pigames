package main

import (
    "bytes"
	"embed"
    "image"
    "image/draw"
    "image/gif"
    "log"
    // "time"

    "periph.io/x/conn/i2c/i2creg"
    "periph.io/x/devices/ssd1306"
    "periph.io/x/devices/ssd1306/image1bit"
    // "periph.io/x/conn/physic"
	"periph.io/x/host"
    

    "github.com/nfnt/resize"
)


var bunny = []byte("GIF87a7\x00@\x00\x80\x01\x00\x00\x00\x00\xff\xff\xff,\x00\x00\x00\x007\x00@\x00\x00\x02\xfe\x84\x8f\xa9\xcb\x16\x1f\u0682\xf4\xc9\xeb*\xc4@k>y\x9bT\x1d\x1ehR\x88\u068cJ\t\xb2\x89\xfbZ\xadL\u04b3]G\x17\x1e\xe2\x05\x194 p\xe7\xcb$;\u0095\x10\x96Z\"\x93P\x9d\xd3ilF\xb1\xb6\xe36Z\xb42\xa921\x18\x1c\x96\x8e\xb9d5\x92\xcd\u02f6G^.\x1c\xfbV\u01b5\xc1\xe5G\x89\xa7\xd7e\xd6\xe3\x87\xe3\x82X\xc6w\xf3\xa4\u0606\xa6\x87\x12y\x068\xe6S\xb5\xf8\x93H78wI\xa8\xb9'\xba6I*iW\xb9\u01b9z:4\xf9gY\xdaJI\n\xcb\x02*y+\xda\xf9\x99yZG\xd49\xdb\xe7\x16\xe3\x9b{\x1c:\u0733\x9c\xda\fl\\\xfc\x9b\x8c\x01\x1d\xdd\xfc\xb5l}\xfdL\u0371\x9d#\u074a\xab\x8cM;->K\x1e\xac\x96\u0337^8\x15\xde\ue747\u02ae\x05\x1exO\\\x98n\u007f\xce\u0554z\xf0\xd0\xe9+\x88,\xc2Aw\xdbtQ[\u05f0\u047fXDRT\xd3\x01\xabO0`Kb\x8e\x9c`\xe4FDH/\x1f\xf1M\xf4Wl\x03\xae\x8c\xa1@\xb1\x9aWR\x90\x8aE/[\x8a\x14\xe1\fO\xc1\x9b8s\xc29iJ\x90@b@}\x86\x839\xaeIG\x97\xf5\x00\x06\xe5X\x0e\x1d\xb3+\x17\x9bf\xbb\xca\xc4\xe60\x8cJi\xb2Lhh\x17O\x84D\u01ce5'\xd3,\x14\xb4ld\xa2-\x00\x00;")


//go:embed images
var image_path embed.FS

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Gray {
    src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
    img := image.NewGray(image.Rect(0, 0, w, h))
    r := src.Bounds()
    r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
    draw.Draw(img, r, src, image.Point{}, draw.Src)
    return img
}

func main() {
    // Load all the drivers:
    if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

 
    // Open a handle to the first available I²C bus:
    bus, err := i2creg.Open("")
 

    // Open a handle to a ssd1306 connected on the I²C bus:
    dev, err := ssd1306.NewI2C(bus, &ssd1306.DefaultOpts)
    
    imgBunnyNRGBA, _ := gif.Decode(bytes.NewReader(bunny))

    imgBunny1bit := image1bit.NewVerticalLSB(imgBunnyNRGBA.Bounds())

    imgBunny1bitLarge := image1bit.NewVerticalLSB(dev.Bounds())

    center := imgBunny1bit.Bounds()
    
    draw.Src.Draw(imgBunny1bitLarge, center.Add(image.Point{X: (128 - center.Dx()) / 2}), imgBunny1bit, image.Point{})
	dev.Draw()
    // imgClear := make([]byte, 128*64/8)

   
    //  _, err = dev.Write(imgClear);
	
    if err != nil {
        log.Fatal(err)
    }
    // g, err := gif.DecodeAll(f)
    // f.Close()
    if err != nil {
        log.Fatal(err)
    }

    
    // Converts every frame to image.Gray and resize them:
    // imgs := make([]*image.Gray, len(g.Image))
    // for i := range g.Image {
    //     imgs[i] = convertAndResizeAndCenter(128, 64, g.Image[i])
    // }

    // Display the frames in a loop:
    // for i := 0; ; i++ {
        // index := i % len(imgs)
        // c := time.After(time.Duration(10*g.Delay[index]) * time.Millisecond)
        // img := imgs[index]
        // dev.Draw(img.Bounds(), img, image.Point{})
        // <-c
    // }
}