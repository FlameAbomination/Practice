package main

import (
	"encoding/json"
	"time"
)

type OrderData struct {
	OrderUID          string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryData `json:"delivery"`
	Payment           PaymentData  `json:"payment"`
	Items             []ItemData   `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	Shardkey          string       `json:"shardkey"`
	SmID              int          `json:"sm_id"`
	DateCreated       time.Time    `json:"date_created"`
	OofShard          string       `json:"oof_shard"`
}

type DeliveryData struct {
	OrderUID string `json:"order_uid,omitempty"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Zip      string `json:"zip"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Region   string `json:"region"`
	Email    string `json:"email"`
}

type PaymentData struct {
	OrderUID     string `json:"order_uid,omitempty"`
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}

type ItemData struct {
	OrderUID    string `json:"order_uid,omitempty"`
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

func JsonToOrder(jsonData []byte) (OrderData, error) {
	var order OrderData
	err := json.Unmarshal(jsonData, &order)
	return order, err
}

func JsonToItem(jsonData []byte) (ItemData, error) {
	var item ItemData
	err := json.Unmarshal(jsonData, &item)
	return item, err
}

func JsonToDelivery(jsonData []byte) (DeliveryData, error) {
	var delivery DeliveryData
	err := json.Unmarshal(jsonData, &delivery)
	return delivery, err
}

func JsonToPayment(jsonData []byte) (PaymentData, error) {
	var payment PaymentData
	err := json.Unmarshal(jsonData, &payment)
	return payment, err
}
