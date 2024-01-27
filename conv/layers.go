package conv

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
)

/** tile size x size */
var layerToSizeMap = map[int][2]int{
	0: {256, 256},
	1: {512, 512},
	2: {512, 1024},
	3: {512, 2048},
	4: {512, 4096},
}

func WriteLayerImage(layerIndex []int, canvases []*image.RGBA, writeDirPath, imgExt string) error {
	if _, err := os.Stat(writeDirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(writeDirPath, os.ModePerm); err != nil {
			return err
		}
	}

	if len(canvases) != faceLen {
		return errors.New("wrong face size")
	}

	for i := range layerIndex {
		level := layerIndex[i]
		fmt.Println("layerIndex: " + strconv.Itoa(level))
		size := layerToSizeMap[level]
		for i, canvas := range canvases {
			if err := write(canvas, level, faceMap[i], size[0], size[1], writeDirPath, imgExt); err != nil {
				return errors.New("Could not write image: " + err.Error())
			}
		}
	}

	return nil
}

func write(img *image.RGBA, level int, tileFace string, tileSize int, size int, writeDirPath, imgExt string) error {

	parts := size / tileSize

	// print sizes
	fmt.Println("tileSize: " + strconv.Itoa(tileSize))
	fmt.Println("size: " + strconv.Itoa(size))
	fmt.Println("parts: " + strconv.Itoa(parts))

	for i := 0; i < parts; i++ {
		for j := 0; j < parts; j++ {
			// id/1/u/0/0.jpg
			targetDir := fmt.Sprintf("%d/%s/%d", level, tileFace[0:1], i)
			fileName := fmt.Sprintf("%d.%s", j, imgExt)

			srcSize := img.Bounds().Dx() / parts

			rect := image.Rect(j*srcSize, i*srcSize, j*srcSize+srcSize, i*srcSize+srcSize)
			sqImg := img.SubImage(rect)

			path := filepath.Join(writeDirPath, targetDir, fileName)

			err := os.MkdirAll(filepath.Join(writeDirPath, targetDir), os.ModePerm)
			if err != nil {
				return errors.New("wrong face size")
			}
			newFile, _ := os.Create(path)

			switch imgExt {
			case "jpg":
				if err := jpeg.Encode(newFile, sqImg, nil); err != nil {
					return err
				}
			case "png":
				if err := png.Encode(newFile, sqImg); err != nil {
					return err
				}
			default:
				return errors.New("Wrong image file format : " + imgExt)
			}
		}
	}

	return nil
}
