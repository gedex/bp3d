package bp3d

import (
	"fmt"
	"reflect"
	"testing"
)

type result struct {
	packed   []*Bin
	unpacked []*Item
}

type testData struct {
	bins        []*Bin
	items       []*Item
	expectation result
}

func TestPack(t *testing.T) {
	// TODO(gedex): move tests data into csv file?
	testCases := []testData{
		// Edge case that needs rotation.
		// from https://github.com/dvdoug/BoxPacker/issues/20
		{
			bins: []*Bin{
				NewBin("Le grande box", 100, 100, 300, 1500),
			},
			items: []*Item{
				NewItem("Item 1", 150, 50, 50, 20),
			},
			expectation: result{
				packed: []*Bin{
					&Bin{
						"Le grande box", 100, 100, 300, 1500,
						[]*Item{
							&Item{"Item 1", 150, 50, 50, 20, RotationType_HDW, Pivot{0, 0, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},

		// test three items fit into smaller bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L12
		{
			bins: []*Bin{
				NewBin("Le petite box", 296, 296, 8, 1000),
				NewBin("Le grande box", 2960, 2960, 80, 10000),
			},
			items: []*Item{
				NewItem("Item 1", 250, 250, 2, 200),
				NewItem("Item 2", 250, 250, 2, 200),
				NewItem("Item 3", 250, 250, 2, 200),
			},
			expectation: result{
				packed: []*Bin{
					&Bin{
						"Le petite box", 296, 296, 8, 1000,
						[]*Item{
							&Item{"Item 1", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 0}},
							&Item{"Item 2", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 2}},
							&Item{"Item 3", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 4}},
						},
					},
					&Bin{
						"Le grande box", 2960, 2960, 80, 10000,
						[]*Item{},
					},
				},
				unpacked: []*Item{},
			},
		},

		// test three items fit into larger bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L36
		{
			bins: []*Bin{
				NewBin("Le petite box", 296, 296, 8, 1000),
				NewBin("Le grande box", 2960, 2960, 80, 10000),
			},
			items: []*Item{
				NewItem("Item 1", 2500, 2500, 20, 2000),
				NewItem("Item 2", 2500, 2500, 20, 2000),
				NewItem("Item 3", 2500, 2500, 20, 2000),
			},
			expectation: result{
				packed: []*Bin{
					&Bin{
						"Le petite box", 296, 296, 8, 1000,
						[]*Item{},
					},
					&Bin{
						"Le grande box", 2960, 2960, 80, 10000,
						[]*Item{
							&Item{"Item 1", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 0}},
							&Item{"Item 2", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 20}},
							&Item{"Item 3", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 40}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},

		// TODO(gedex): five items packed into two large bins and one small bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L60

		// 1 bin with 7 items fit into.
		// from https://github.com/bom-d-van/binpacking/blob/master/binpacking_test.go
		{
			bins: []*Bin{
				NewBin("Bin 1", 220, 160, 100, 110),
			},
			items: []*Item{
				NewItem("Item 1", 20, 100, 30, 10),
				NewItem("Item 2", 100, 20, 30, 10),
				NewItem("Item 3", 20, 100, 30, 10),
				NewItem("Item 4", 100, 20, 30, 10),
				NewItem("Item 5", 100, 20, 30, 10),
				NewItem("Item 6", 100, 100, 30, 10),
				NewItem("Item 7", 100, 100, 30, 10),
			},
			expectation: result{
				packed: []*Bin{
					&Bin{
						"Bin 1", 220, 160, 100, 110,
						[]*Item{
							&Item{"Item 7", 100, 100, 30, 10, RotationType_WHD, Pivot{0, 0, 0}},
							&Item{"Item 6", 100, 100, 30, 10, RotationType_WHD, Pivot{100, 0, 0}},
							&Item{"Item 2", 100, 20, 30, 10, RotationType_HWD, Pivot{200, 0, 0}},
							&Item{"Item 3", 20, 100, 30, 10, RotationType_HWD, Pivot{0, 100, 0}},
							&Item{"Item 4", 100, 20, 30, 10, RotationType_WHD, Pivot{100, 100, 0}},
							&Item{"Item 5", 100, 20, 30, 10, RotationType_HDW, Pivot{200, 100, 0}},
							&Item{"Item 1", 20, 100, 30, 10, RotationType_HWD, Pivot{100, 120, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},
	}

	for _, tc := range testCases {
		testPack(t, tc)
	}
}

func testPack(t *testing.T, td testData) {
	packer := NewPacker()
	packer.AddBin(td.bins...)
	packer.AddItem(td.items...)

	if err := packer.Pack(); err != nil {
		t.Fatalf("Got error: %v", err)
	}

	if !reflect.DeepEqual(packer.Bins, td.expectation.packed) {
		t.Errorf("\nGot:\n%+v\nwant:\n%+v", formatBins(packer.Bins, packer.UnfitItems), formatBins(td.expectation.packed, td.expectation.unpacked))
	}
}

func formatBins(bins []*Bin, unpacked []*Item) string {
	var s string
	for _, b := range bins {
		s += fmt.Sprintln(b)
		s += fmt.Sprintln(" packed items:")
		for _, i := range b.Items {
			s += fmt.Sprintln("  ", i)
		}

		s += fmt.Sprintln(" unpacked items:")
		for _, i := range unpacked {
			s += fmt.Sprintln("  ", i)
		}
	}
	return s
}
