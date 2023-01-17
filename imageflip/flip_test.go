package imageflip

import (
    "testing"
	"image/gif"
	"image/color"
	"os"
	"reflect"
)

var b = color.Black
var w = color.White
	// expected output
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
	// expected output
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

	img := imageToGrid(avatarGif.Image[0])

    if !reflect.DeepEqual(img, output)  || err != nil {
        t.Errorf(`imageToGrid(avatarGif.Image[0]) = %q, and not %q`, img,  output)
    }
}

func TestFlip(t *testing.T) {
	

	avatarpath := "../images/avatar.gif"
	f, err := os.Open(avatarpath)
	avatarGif, err := gif.DecodeAll(f)
	f.Close()

	resultImg := Flip(avatarGif.Image[0])

	result := imageToGrid(resultImg) // the flip returns an image but it's easier to compare using grids.

    if !reflect.DeepEqual(result, flipOuput)  || err != nil {
        t.Errorf(`Flip(avatarGif.Image[0]) = %q, and not %q`, result,  flipOuput)
    }
}