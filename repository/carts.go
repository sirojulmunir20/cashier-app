package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
)

type CartRepository struct {
	db db.DB
}

func NewCartRepository(db db.DB) CartRepository {
	return CartRepository{db}
}

func (u *CartRepository) ReadCart() ([]model.Cart, error) {
	records, err := u.db.Load("carts")
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("Cart not found!")
	}

	var cart []model.Cart
	err = json.Unmarshal([]byte(records), &cart)
	if err != nil {
		return nil, err
	}

	return cart, nil
}

func (u *CartRepository) UpdateCart(cart model.Cart) error {
	readCart, err := u.ReadCart()
	if err != nil {
		return err
	}
	
	for i := 0; i < len(readCart); i++ {
		if readCart[i].Name == cart.Name {
			readCart[i] = cart
		}
	}
	jsonData, err := json.Marshal(readCart)
	if err != nil {
		return err
	}

	err = u.db.Save("carts", jsonData)
	if err != nil {
		return err
	}

	return nil
}

func (u *CartRepository) AddCart(cart model.Cart) error {
	data, _ := u.ReadCart()
	data = append(data, cart)
	jsonCart, _ := json.Marshal(data)
	_= u.db.Save("carts", jsonCart)
	return nil // TODO: replace this
}

func (u *CartRepository) ResetCarts() error {
	err := u.db.Reset("carts", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *CartRepository) CartUserExist(name string) (model.Cart, error) {
	listcCart, err := u.ReadCart()
	if err != nil {
		return model.Cart{}, err
	}
	for _, element := range listcCart {
		if element.Name == name {
			return element, nil
		}
	}
	return model.Cart{}, fmt.Errorf("Cart Empty!")
}
