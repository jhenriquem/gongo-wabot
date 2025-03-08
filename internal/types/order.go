package types

type client struct {
	Name        string
	PhoneNumber string
}

type Order struct {
	Client        client
	Book          Book
	Address       string
	PaymentMethod string
	Data          string
	Price         float64
}
