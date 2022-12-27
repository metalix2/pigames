package main

import (
	"embed"
    "image"
    "image/draw"
    "image/gif"
    "image/color/palette"
    "image/color"
    "log"

    "time"

    "periph.io/x/conn/i2c/i2creg"
    "periph.io/x/devices/ssd1306"
    // "periph.io/x/devices/ssd1306/image1bit"
	"periph.io/x/host"
    "periph.io/x/conn/gpio"
	"periph.io/x/conn/gpio/gpioreg"
    

    "github.com/nfnt/resize"
)


// var bunny = []byte("GIF87a7\x00@\x00\x80\x01\x00\x00\x00\x00\xff\xff\xff,\x00\x00\x00\x007\x00@\x00\x00\x02\xfe\x84\x8f\xa9\xcb\x16\x1f\u0682\xf4\xc9\xeb*\xc4@k>y\x9bT\x1d\x1ehR\x88\u068cJ\t\xb2\x89\xfbZ\xadL\u04b3]G\x17\x1e\xe2\x05\x194 p\xe7\xcb$;\u0095\x10\x96Z\"\x93P\x9d\xd3ilF\xb1\xb6\xe36Z\xb42\xa921\x18\x1c\x96\x8e\xb9d5\x92\xcd\u02f6G^.\x1c\xfbV\u01b5\xc1\xe5G\x89\xa7\xd7e\xd6\xe3\x87\xe3\x82X\xc6w\xf3\xa4\u0606\xa6\x87\x12y\x068\xe6S\xb5\xf8\x93H78wI\xa8\xb9'\xba6I*iW\xb9\u01b9z:4\xf9gY\xdaJI\n\xcb\x02*y+\xda\xf9\x99yZG\xd49\xdb\xe7\x16\xe3\x9b{\x1c:\u0733\x9c\xda\fl\\\xfc\x9b\x8c\x01\x1d\xdd\xfc\xb5l}\xfdL\u0371\x9d#\u074a\xab\x8cM;->K\x1e\xac\x96\u0337^8\x15\xde\ue747\u02ae\x05\x1exO\\\x98n\u007f\xce\u0554z\xf0\xd0\xe9+\x88,\xc2Aw\xdbtQ[\u05f0\u047fXDRT\xd3\x01\xabO0`Kb\x8e\x9c`\xe4FDH/\x1f\xf1M\xf4Wl\x03\xae\x8c\xa1@\xb1\x9aWR\x90\x8aE/[\x8a\x14\xe1\fO\xc1\x9b8s\xc29iJ\x90@b@}\x86\x839\xaeIG\x97\xf5\x00\x06\xe5X\x0e\x1d\xb3+\x17\x9bf\xbb\xca\xc4\xe60\x8cJi\xb2Lhh\x17O\x84D\u01ce5'\xd3,\x14\xb4ld\xa2-\x00\x00;")

//go:embed images
var image_path embed.FS

// convertAndResizeAndCenter takes an image, resizes and centers it on a
// image.Gray of size w*h.
func convertAndResizeAndCenter(w, h int, src image.Image) *image.Paletted {
    src = resize.Thumbnail(uint(w), uint(h), src, resize.Bicubic)
    r := src.Bounds()
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)
    // img := image.NewGray(image.Rect(0, 0, w, h))
    
    r = r.Add(image.Point{(w - r.Max.X) / 2, (h - r.Max.Y) / 2})
    draw.Draw(img, r, src, image.Point{}, draw.Src)
    // draw.Draw(palettedImage, palettedImage.Rect, simage, bounds.Min, draw.Over)
    return img
}

func flip(img image.Image) image.Image {
    // create and return a grid of pixels
    var grid  [][]color.Color
    size := img.Bounds().Size()
    for i := 0; i < size.X; i++ {
        var y []color.Color
        for j := 0; j < size.Y; j++ {
            y = append(y, img.At(i, j))
        }
        grid = append(grid, y)
    }

    // Flips Verticaly
    lenX := len(grid)-1
    for x := 0; x <= len(grid) / 2; x++ {
        col := grid[x]
        for  y := 0; y < len(col); y++ {
            var t = grid[x][y]
            grid[x][y] = grid[lenX-x][y]
            grid[lenX-x][y] = t
        }
    }
  
    xlen, ylen := len(grid), len(grid[0])
    rect := image.Rect(0, 0, xlen, ylen)
    img2 := image.NewPaletted(rect, palette.Plan9)
    for x := 0; x < xlen; x++ {
        for y := 0; y < ylen; y++ {
            img2.Set(x, y, grid[x][y])
        }
    }
    return img2.SubImage(rect)
}

func drawCanvas(w, h int, src image.Image, pos_x int, pos_y int, dir int) *image.Paletted {
    r := src.Bounds()
    img := image.NewPaletted(image.Rect(0, 0, w, h), palette.Plan9)
    // img := image.NewGray(image.Rect(0, 0, w, h))
    
    r = r.Add(image.Point{(0 - pos_x), (0 - pos_y)})
    if dir > 0 {
        draw.Draw(img, r, flip(src), image.Point{}, draw.Src) 
    } else {
        draw.Draw(img, r, src, image.Point{}, draw.Src)
    }
    // draw.Draw(palettedImage, palettedImage.Rect, simage, bounds.Min, draw.Over)
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

    // Lookup a pin by its number:
	p_r := gpioreg.ByName("GPIO23")
    p_l := gpioreg.ByName("GPIO27")
    p_u := gpioreg.ByName("GPIO17")
    p_d := gpioreg.ByName("GPIO22")
	if p_r == nil {
		log.Fatal("Failed to find GPIO27")
	}

	log.Println("%s: %s\n", p_r, p_r.Function())

	// Set it as input, with an internal pull up resistor:
	if err := p_r.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_l.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_u.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    if err := p_d.In(gpio.PullUp, gpio.BothEdges); err != nil {
		log.Fatal(err)
	}
    // Display the frames in a loop:
    var xx int = 0
    var yy int = 0
    var dir = 0
    var fps = 10
    ts := time.Now();
    for i := 1; ;  {
        index := i % len(g.Image)
        c := time.After(time.Duration(10*fps) * time.Millisecond)
        if ts.Add(time.Duration(10 * g.Delay[index]) * time.Millisecond).Sub(time.Now()) < time.Duration(10 * 1) * time.Millisecond {
            i++
            ts = time.Now();
            index = i % len(g.Image)
        }
        img := drawCanvas(128, 64, g.Image[index], xx, yy, dir)
        if p_r.Read() ==  gpio.Low {
            log.Println(img.Bounds());
            dir = 1
            xx -= 2
        }
        if p_l.Read() ==  gpio.Low {
            log.Println(img.Bounds());
            dir = 0
            xx += 2
        }
        if p_u.Read() ==  gpio.Low {
            log.Println(img.Bounds());
            yy += 2
        }
        if p_d.Read() ==  gpio.Low {
            log.Println(img.Bounds());
            yy -= 2
        }
        
        dev.Draw(img.Bounds(), img, image.Point{})
        // if i == len(g.Image) + 1 {
        //     break
        // }
        <-c
    }
}