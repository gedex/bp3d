package bp3d

import (
	"fmt"
	"reflect"
	"sort"
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

	}

	for _, tc := range testCases {
		testPack(t, tc)
	}
}

func TestPackBomDVanFixture(t *testing.T) {
	packer := NewPacker()
	packer.AddBin(NewBin("Bin 1", 220, 160, 100, 110))
	packer.AddItem(
		NewItem("Item 1", 20, 100, 30, 10),
		NewItem("Item 2", 100, 20, 30, 10),
		NewItem("Item 3", 20, 100, 30, 10),
		NewItem("Item 4", 100, 20, 30, 10),
		NewItem("Item 5", 100, 20, 30, 10),
		NewItem("Item 6", 100, 100, 30, 10),
		NewItem("Item 7", 100, 100, 30, 10),
	)

	if err := packer.Pack(); err != nil {
		t.Fatalf("Got error: %v", err)
	}

	assertPackedItemNames(t, packer.Bins, packer.UnfitItems, []string{
		"Item 1",
		"Item 2",
		"Item 3",
		"Item 4",
		"Item 5",
		"Item 6",
		"Item 7",
	})
	assertNoOutOfBoundsItems(t, packer.Bins)
}

func TestPackPackriftEcommerceCartonFixture(t *testing.T) {
	packer := NewPacker()
	packer.AddBin(NewBin("Packrift 40x20x20 carton", 40, 20, 20, 75))

	// Fixture source:
	// https://packrift.github.io/packaging-optimization-benchmark-corpus/cartonization-solver-fixtures.html
	// Dimensions are in inches and weights are in pounds.
	packer.AddItem(
		NewItem("demo-small-item", 7.5, 4.5, 3.5, 1),
		NewItem("demo-flat-item-1", 15, 7, 2.5, 1.5),
		NewItem("demo-flat-item-2", 15, 7, 2.5, 1.5),
		NewItem("demo-long-item", 21, 8.5, 5, 3),
		NewItem("demo-bulk-item-1", 18, 12, 5.5, 4),
		NewItem("demo-bulk-item-2", 18, 12, 5.5, 4),
		NewItem("demo-bulk-item-3", 18, 12, 5.5, 4),
		NewItem("demo-bulk-item-4", 18, 12, 5.5, 4),
	)

	if err := packer.Pack(); err != nil {
		t.Fatalf("Got error: %v", err)
	}

	assertPackedItemNames(t, packer.Bins, packer.UnfitItems, []string{
		"demo-small-item",
		"demo-flat-item-1",
		"demo-flat-item-2",
		"demo-long-item",
		"demo-bulk-item-1",
		"demo-bulk-item-2",
		"demo-bulk-item-3",
		"demo-bulk-item-4",
	})
	assertNoOutOfBoundsItems(t, packer.Bins)

	if len(packer.Bins[0].Items) != 8 {
		t.Fatalf("Packrift fixture packed %d items, want 8", len(packer.Bins[0].Items))
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

func assertPackedItemNames(t *testing.T, bins []*Bin, unpacked []*Item, want []string) {
	t.Helper()
	if len(unpacked) != 0 {
		t.Fatalf("Got unpacked items: %v", unpacked)
	}

	got := make([]string, 0, len(want))
	for _, b := range bins {
		for _, i := range b.Items {
			got = append(got, i.Name)
		}
	}
	sort.Strings(got)
	sort.Strings(want)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Got packed items %v, want %v", got, want)
	}
}

func assertNoOutOfBoundsItems(t *testing.T, bins []*Bin) {
	t.Helper()

	for _, b := range bins {
		for _, i := range b.Items {
			d := i.GetDimension()
			if i.Position[0] < 0 || i.Position[1] < 0 || i.Position[2] < 0 {
				t.Fatalf("%s has negative position %v in %s", i.Name, i.Position, b.Name)
			}
			if i.Position[0]+d[0] > b.Width || i.Position[1]+d[1] > b.Height || i.Position[2]+d[2] > b.Depth {
				t.Fatalf("%s at %v with dimension %v is outside %s", i.Name, i.Position, d, b)
			}
		}
	}
}
