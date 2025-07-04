package main

import (
	"fmt"
	"testing"
)

func TestGenerateNewImage(t *testing.T) {
	imagePath := "https://img-va.myshopline.com/image/store/1699869479480/Frog-Family-Wooden-Jigsaw-Puzzle.jpeg?w=1080&h=1080"
	ctid, err := GenerateNewImage(imagePath)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	fmt.Println("GenerateNewImage return image ctid", ctid)
}

func TestGetNewImage(t *testing.T) {
	var ctid int64 = 266597818931087
	url, err := GetNewImageInfo(ctid)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}
	fmt.Println("GetNewImageInfo return image url", url)
}
