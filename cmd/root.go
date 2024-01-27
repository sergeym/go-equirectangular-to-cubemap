package cmd

import (
	"fmt"
	"github.com/blackironj/panorama/conv"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const defaultEdgeLen = 1024

var (
	inFilePath string
	outFileDir string
	edgeLen    int

	rootCmd = &cobra.Command{
		Use:   "panorama",
		Short: "convert equirectangular panorama img to Cubemap img",
		Run: func(cmd *cobra.Command, args []string) {
			if inFilePath == "" {
				er("Need a image for converting")
			}

			// check for inFilePath is directory or file
			// if directory, convert all files in the directory
			// if file, convert the file

			isDir, err := isDirectory(inFilePath)
			if err != nil {
				er(err)
			}

			var files []string

			if isDir {
				fmt.Println(fmt.Sprintf("`%s` is Directory", inFilePath))

				// get list of files in the directory
				inDirFiles, err := os.ReadDir(inFilePath)
				if err != nil {
					er(err)
				}

				for _, entry := range inDirFiles {
					files = append(files, entry.Name())
				}
			} else {
				fmt.Println("File")
				files = append(files, inFilePath)
			}

			for _, filePath := range files {

				fmt.Println("Read a image...")
				inImage, ext, err := conv.ReadImage(filepath.Join(inFilePath, filePath))
				if err != nil {
					er(err)
				}

				s := spinner.New(spinner.CharSets[33], 100*time.Millisecond)
				s.FinalMSG = "Complete converting!\n"
				s.Prefix = "Converting..."

				s.Start()
				canvases := conv.ConvertEquirectangularToCubeMap(edgeLen, inImage)

				fmt.Println("Write images...")

				targetDir := filepath.Join(outFileDir, strings.TrimSuffix(filePath, path.Ext(filePath)))

				if err := conv.WriteImage(canvases, targetDir, ext); err != nil {
					er(err)
				}

				fmt.Println("Creating preview...")
				preview, err := conv.GeneratePreview(canvases)
				if err != nil {
					er(err)
				}

				fmt.Println("Writing preview...")
				if err := conv.WritePreviewImage(preview, targetDir, ext); err != nil {
					er(err)
				}

				fmt.Println("Writing layers...")
				reqListOfIndex := []int{1, 2, 3, 4}
				if err := conv.WriteLayerImage(reqListOfIndex, canvases, targetDir, ext); err != nil {
					er(err)
				}

				s.Stop()
			}
		},
	}
)

func init() {
	rootCmd.Flags().StringVarP(&inFilePath, "in", "i", "", "in image file path (required)")
	rootCmd.Flags().StringVarP(&outFileDir, "out", "o", ".", "out file dir path")
	rootCmd.Flags().IntVarP(&edgeLen, "len", "l", defaultEdgeLen, "edge length of a cube face")
}

func er(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		er(err)
	}
}

func isDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}
