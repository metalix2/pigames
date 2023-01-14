package imageFlip

import (
    "testing"
	"image/gif"
	"image/color"
	"os"
	"reflect"
)

var b = color.Black
var w = color.White
var output = [][]color.Color{
	{b, b, b, w, w, w, w, w, b, b, b, b, b},
	{b, b, w, w, b, b, b, w, w, b, b, b, b},
	{b, w, w, b, w, b, w, b, b, b, b, b, w},
	{b, w, w, b, w, w, w, b, w, w, b, w, w},
	{b, w, w, b, w, b, w, b, w, w, b, b, b},
	{b, w, w, b, w, w, w, b, w, w, b, w, w},
	{b, w, w, w, b, b, w, b, b, w, b, b, w},
	{b, w, w, w, w, w, b, b, b, b, b, b, b},
	{w, b, w, w, w, w, w, w, b, b, b, b, b},
	{b, b, w, b, w, w, w, b, w, b, b, b, b},
	{b, b, b, b, b, b, w, b, b, b, b, b, b},
}

var flipOuput = [][]color.Color{
	{b, b, b, b, b, b, w, b, b, b, b, b, b},
	{b, b, w, b, w, w, w, b, w, b, b, b, b},
	{w, b, w, w, w, w, w, w, b, b, b, b, b},
	{b, w, w, w, w, w, b, b, b, b, b, b, b},
	{b, w, w, w, b, b, w, b, b, w, b, b, w},
	{b, w, w, b, w, w, w, b, w, w, b, w, w},
	{b, w, w, b, w, b, w, b, w, w, b, b, b},
	{b, w, w, b, w, w, w, b, w, w, b, w, w},
	{b, w, w, b, w, b, w, b, b, b, b, b, w},
	{b, b, w, w, b, b, b, w, w, b, b, b, b},
	{b, b, b, w, w, w, w, w, b, b, b, b, b},
}

func TestImageToGrid(t *testing.T) {
	
	avatarpath := "../images/avatar.gif"
	f, err := os.Open(avatarpath)
	avatarGif, err := gif.DecodeAll(f)
	f.Close()

	// expected output

	img := imageToGrid(avatarGif.Image[0])

    if !reflect.DeepEqual(img, output)  || err != nil {
        t.Errorf(`imageToGrid(avatarGif.Image[0]) = %q, %q`, img,  output)
    }
}

func TestImageFlip(t *testing.T) {
	

	avatarpath := "../images/avatar.gif"
	f, err := os.Open(avatarpath)
	avatarGif, err := gif.DecodeAll(f)
	f.Close()

	// expected output

	resultImg := flip(avatarGif.Image[0])

	result := imageToGrid(resultImg)

    if !reflect.DeepEqual(result, flipOuput)  || err != nil {
        t.Errorf(`imageToGrid(avatarGif.Image[0]) = %q, %q`, result,  flipOuput)
    }
}