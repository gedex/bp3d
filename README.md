bp3d
====

3D Bin Packing implementation based on [this paper](http://www.cs.ukzn.ac.za/publications/erick_dube_507-034.pdf). The code is based on [binpacking by bom-d-van](https://github.com/bom-d-van/binpacking) but
modified to allow flexible bins and use `float64` instead of `int`.

## Install

```
go get github.com/Automattic/bp3d
```

## Usage

```
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

// Each bin, b, in p.Bins might have packed items in b.Items
```

See [`example/example.go`](./example/example.go)

## Credit

* http://www.cs.ukzn.ac.za/publications/erick_dube_507-034.pdf
* https://github.com/bom-d-van/binpacking

## License

[MIT](./LICENSE)
