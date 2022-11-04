package main

import (
	"fmt"
	"log"

	"github.com/spicescode/bp3d"
)

func main() {
	p := bp3d.NewPacker()

	// Add bins.
	p.AddBin(bp3d.NewBin("Small Bin", 10, 15, 20, 100))
	p.AddBin(bp3d.NewBin("Medium Bin", 100, 150, 200, 1000))

	// Add items.
	p.AddItem(bp3d.NewItem("Item 1", 2, 2, 1, 2))
	p.AddItem(bp3d.NewItem("Item 2", 3, 3, 2, 3))

	// Pack items to bins.
	if err := p.Pack(); err != nil {
		log.Fatal(err)
	}

	// Will output:
	//
	// Small Bin(10x15x20, max_weight:100)
	//  packed items:
	//    Item 2(3x3x2, weight: 3) pos(0,0,0) rt(RotationType_WHD (w,h,d))
	//    Item 1(2x2x1, weight: 2) pos(3,0,0) rt(RotationType_WHD (w,h,d))
	// Medium Bin(100x150x200, max_weight:1000)
	//  packed items:
	displayPacked(p.Bins)
}

func displayPacked(bins []*bp3d.Bin) {
	for _, b := range bins {
		fmt.Println(b)
		fmt.Println(" packed items:")
		for _, i := range b.Items {
			fmt.Println("  ", i)
		}
	}
}
