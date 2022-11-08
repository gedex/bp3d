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
	name        string
	bins        []*Bin
	items       []*Item
	errExpected bool
	expectation result
}

func TestPack(t *testing.T) {
	// TODO(gedex): move tests data into csv file?
	testCases := []testData{
		// Edge case that needs rotation.
		// from https://github.com/dvdoug/BoxPacker/issues/20
		{
			name: "Edge case that needs rotation.",
			bins: []*Bin{
				NewBin("Le grande box", 100, 100, 300, 1500),
			},
			items: []*Item{
				NewItem("Item 1", 150, 50, 50, 20),
			},
			expectation: result{
				packed: []*Bin{
					{
						"Le grande box", 100, 100, 300, 1500,
						[]*Item{
							{"Item 1", 150, 50, 50, 20, RotationType_HDW, Pivot{0, 0, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},

		// test three items fit into smaller bin.
		// from https://github.com/dvdoug/BoxPacker/blob/master/tests/PackerTest.php#L12
		{
			name: "test three items fit into smaller bin.",
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
					{
						"Le petite box", 296, 296, 8, 1000,
						[]*Item{
							{"Item 1", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 0}},
							{"Item 2", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 2}},
							{"Item 3", 250, 250, 2, 200, RotationType_WHD, Pivot{0, 0, 4}},
						},
					},
					{
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
			name: "test three items fit into larger bin.",
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
					{
						"Le petite box", 296, 296, 8, 1000,
						[]*Item{},
					},
					{
						"Le grande box", 2960, 2960, 80, 10000,
						[]*Item{
							{"Item 1", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 0}},
							{"Item 2", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 20}},
							{"Item 3", 2500, 2500, 20, 2000, RotationType_WHD, Pivot{0, 0, 40}},
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
			name: "1 bin with 7 items fit into.",
			bins: []*Bin{
				NewBin("Bin 1", 220, 160, 100, 110),
			},
			items: []*Item{
				NewItem("Item 7", 100, 100, 30, 10),
				NewItem("Item 6", 100, 100, 30, 10),
				NewItem("Item 2", 100, 20, 30, 10),
				NewItem("Item 3", 20, 100, 30, 10),
				NewItem("Item 4", 100, 20, 30, 10),
				NewItem("Item 5", 100, 20, 30, 10),
				NewItem("Item 1", 20, 100, 30, 10),
			},
			expectation: result{
				packed: []*Bin{
					{
						"Bin 1", 220, 160, 100, 110,
						[]*Item{
							{"Item 7", 100, 100, 30, 10, RotationType_WHD, Pivot{0, 0, 0}},
							{"Item 6", 100, 100, 30, 10, RotationType_WHD, Pivot{100, 0, 0}},
							{"Item 2", 100, 20, 30, 10, RotationType_HWD, Pivot{200, 0, 0}},
							{"Item 3", 20, 100, 30, 10, RotationType_HWD, Pivot{0, 100, 0}},
							{"Item 4", 100, 20, 30, 10, RotationType_WHD, Pivot{100, 100, 0}},
							{"Item 5", 100, 20, 30, 10, RotationType_HDW, Pivot{200, 100, 0}},
							{"Item 1", 20, 100, 30, 10, RotationType_HWD, Pivot{100, 120, 0}},
						},
					},
				},
				unpacked: []*Item{},
			},
		},
		// invalid volume error
		{
			name: "invalid volume error",
			bins: []*Bin{
				NewBin("Bin 1", 220, 160, 100, 110),
			},
			items: []*Item{
				NewItem("Item 1", 230, 160, 110, 10),
				NewItem("Item 2", 230, 100, 30, 10),
			},
			errExpected: true,
		},
		// Unfit error
		{
			name: "Unfit error",
			bins: []*Bin{
				NewBin("Bin 1", 220, 160, 100, 110),
			},
			items: []*Item{
				NewItem("Item 1", 230, 100, 30, 10),
				NewItem("Item 2", 230, 100, 30, 10),
				NewItem("Item 3", 230, 100, 30, 10),
				NewItem("Item 4", 230, 100, 30, 10),
			},
			errExpected: true,
		},
		{
			name: "large unfit",
			bins: []*Bin{
				NewBin("120size-0", 53.5, 24.0, 38.5, 20.0),
				NewBin("120size-1", 53.5, 24.0, 38.5, 20.0),
				NewBin("120size-2", 53.5, 24.0, 38.5, 20.0),
				NewBin("60size-0", 26.4, 11.4, 19.4, 5.0),
			},
			items: []*Item{
				NewItem("Item 1", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 2", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 3", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 4", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 5", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 6", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 7", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 8", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 9", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 10", 10.0, 10.0, 10.0, 2.0),
				NewItem("Item 11", 50.0, 30.0, 10.0, 2.0),
				NewItem("Item 12", 50.0, 30.0, 10.0, 2.0),
				NewItem("Item 13", 50.0, 30.0, 10.0, 2.0),
			},
			errExpected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			testPack(t, tc)
		})
	}
}

func testPack(t *testing.T, td testData) {
	packer := NewPacker()
	packer.AddBin(td.bins...)
	packer.AddItem(td.items...)

	if err := packer.Pack(); err != nil {
		if td.errExpected {
			return
		}
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
