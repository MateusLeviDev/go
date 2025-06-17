package main

type Inventory map[string]int

type InventoryErr string

func (e InventoryErr) Error() string {
	return string(e)
}

const (
	ErrProductNotFound = InventoryErr("product not found")
	ErrProductExists   = InventoryErr("product already exists")
)

func (inv Inventory) Get(name string) (int, error) {
	qty, ok := inv[name]

	if !ok {
		return 0, ErrProductNotFound
	}
	return qty, nil
}
