package conv

import (
	"errors"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
)

var previewToFaceMap = map[int]int{
	0: 0,
	1: 5,
	2: 2,
	3: 1,
	4: 3,
	5: 4,
}

func GeneratePreview(canvases []*image.RGBA) (*image.RGBA, error) {
	if len(canvases) != faceLen {
		return nil, errors.New("wrong face size")
	}

	preview := image.NewRGBA(image.Rect(0, 0, 256, 1536))

	for i := 0; i < faceLen; i++ {
		offset := image.Pt(0, i*256)
		newImage := resize.Resize(256, 0, canvases[previewToFaceMap[i]], resize.Lanczos3)
		draw.Draw(preview, newImage.Bounds().Add(offset), newImage, image.ZP, draw.Over)
	}

	return preview, nil
}

func WritePreviewImage(preview *image.RGBA, writeDirPath, imgExt string) error {
	if _, err := os.Stat(writeDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(writeDirPath, os.ModePerm); err != nil {
			return err
		}
	}

	path := filepath.Join(writeDirPath, "preview."+imgExt)
	newFile, _ := os.Create(path)

	switch imgExt {
	case "jpg":
		if err := jpeg.Encode(newFile, preview, nil); err != nil {
			return err
		}
	case "png":
		if err := png.Encode(newFile, preview); err != nil {
			return err
		}
	default:
		return errors.New("Wrong image file format : " + imgExt)
	}

	return nil
}
