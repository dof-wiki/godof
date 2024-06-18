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
	os.RemoveAll("output")
	os.Mkdir("output", os.ModePerm)
	//path := "/Users/ziipin/Downloads/新版婚纱皮肤导入/！！！newmarryme.NPK"
	//path := "D:\\Games\\dnf\\DNF\\ImagePacks2\\!!!!+(登入).NPK"
	//path := "D:\\Games\\dnf\\DNF\\ImagePacks2\\!!!双天空城_旧版.NPK"
	path := "/Users/ziipin/Downloads/sprite_character_fighter_effect_poisonexplosioncustom.NPK"
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	n, err := npk.Open(f)
	if err != nil {
		panic(err)
	}
	for j, file := range n.Files {
		img, err := file.ToIMG()
		if err != nil {
			continue
		}
		for i, item := range img.Images {
			image, err := img.Build(item)
			if err != nil {
				panic(err)
			}
			saveImage(image, fmt.Sprintf("output/%d_%d.png", j, i))
		}
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
