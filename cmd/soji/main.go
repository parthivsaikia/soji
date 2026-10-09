package main

import (
	"fmt"
	"log"

	"github.com/parthivsaikia/soji/internal/discovery/kubernetes"
)

func main() {
	fmt.Println("finding images...")
	images, err := kubernetes.FindCurrentlyUsedImage()
	if err != nil {
		log.Printf("error in finding image: %v", err)
	}
	fmt.Println("found images...")
	for i, image := range images {
		fmt.Printf("%d image: %s\n", i+1, image)
	}
}
