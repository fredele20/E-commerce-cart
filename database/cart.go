package database

import "errors"

var (
	ErrProductNotFound       = errors.New("product not found")
	ErrProductDecodingFailed = errors.New("product decoding failed")
	ErrUserIdNotValid        = errors.New("this user is not valid")
	ErrCantRemoveCartItem    = errors.New("cannot remove this item from the cart")
	ErrCantGetItem           = errors.New("unable to get cart item")
	ErrCantBuyCartItem       = errors.New("cannot update the purchase")
	ErrCantUpdateUser        = errors.New("cannot add this product to the cart")
)

func AddProductToCart() {}

func RemoveCartItem() {}

func BuyItemFromCart() {}

func InstantBuyer() {}
