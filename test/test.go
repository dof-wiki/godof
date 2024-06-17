package main

import (
	"bufio"
	"fmt"
	"github.com/dof-wiki/godof/npk"
	image2 "image"
	"image/png"
	"os"
)

func main() {
	path := "/Users/ziipin/Downloads/新版婚纱皮肤导入/！！！newmarryme.NPK"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := npk.Open(f)
	if err != nil {
		panic(err)
	}
	for _, file := range n.Files {
		img, err := file.ToIMG()
		if err != nil {
			panic(err)
		}
		for i, item := range img.Images {
			image, err := img.Build(item)
			if err != nil {
				panic(err)
			}
			saveImage(image, fmt.Sprintf("output/%d.png", i))
		}
		break
	}
}

func saveImage(image image2.Image, name string) {
	outFile, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	b := bufio.NewWriter(outFile)
	err = png.Encode(b, image)
	if err != nil {
		panic(err)
	}
	err = b.Flush()
	if err != nil {
		panic(err)
	}
}
