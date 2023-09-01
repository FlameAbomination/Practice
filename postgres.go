package main

import (
	"context"
	"sync"

	pgx "github.com/jackc/pgx/v5"
)

// Кэш и mutex для защиты от параллельной записи в кэш
var mutex sync.Mutex
var cache map[string]OrderData

// Чтение данных из кэша
func GetOrder(uid string) (OrderData, bool) {
	mutex.Lock()
	defer mutex.Unlock()
	val, ok := cache[uid]
	if ok {
		return val, true
	}
	return val, false
}

// Чтение всех таблиц из базы и инициализация кэша
func LoadDatabase(conn *pgx.Conn) error {
	mutex.Lock()
	cache = make(map[string]OrderData)
	items, err := GetAllItems(conn)
	if err != nil {
		return err
	}
	deliveries, err := GetAllDelivery(conn)
	if err != nil {
		return err
	}
	payments, err := GetAllPayments(conn)
	if err != nil {
		return err
	}
	orders, err := GetAllOrders(conn)
	if err != nil {
		return err
	}
	for _, order := range orders {
		cache[order.OrderUID] = order
	}

	for _, item := range items {
		// OrderUID необходим только для связывания заказа с товаром и т.д.
		// поэтому его можно удалить перед чтением в кэш
		order := cache[item.OrderUID]
		item.OrderUID = ""
		order.Items = append(order.Items, item)
		cache[order.OrderUID] = order
	}

	for _, delivery := range deliveries {
		order := cache[delivery.OrderUID]
		delivery.OrderUID = ""
		order.Delivery = delivery
		cache[order.OrderUID] = order
	}

	for _, payment := range payments {
		order := cache[payment.OrderUID]
		payment.OrderUID = ""
		order.Payment = payment
		cache[order.OrderUID] = order
	}

	mutex.Unlock()
	return nil
}

// TODO: изучить возможность использования generic'ов
// Данные передаются в виде json для унификации входных данных,
// так как NATS тоже передаёт данные в json
func GetAllItems(conn *pgx.Conn) ([]ItemData, error) {
	var rowSlice []ItemData
	rows, err := conn.Query(context.Background(),
		`SELECT row_to_json(row) FROM items row`)
	if err != nil {
		return rowSlice, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		rows.Scan(&data)
		item, err := JsonToItem(data)
		if err != nil {
			return rowSlice, err
		}
		rowSlice = append(rowSlice, item)
	}
	if err := rows.Err(); err != nil {
		return rowSlice, err
	}
	return rowSlice, nil
}

func GetAllDelivery(conn *pgx.Conn) ([]DeliveryData, error) {
	var rowSlice []DeliveryData
	rows, err := conn.Query(context.Background(),
		`SELECT row_to_json(row) FROM delivery row`)

	if err != nil {
		return rowSlice, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		rows.Scan(&data)
		delivery, err := JsonToDelivery(data)
		if err != nil {
			return rowSlice, err
		}

		rowSlice = append(rowSlice, delivery)
	}
	if err := rows.Err(); err != nil {
		return rowSlice, err
	}
	return rowSlice, nil
}

func GetAllPayments(conn *pgx.Conn) ([]PaymentData, error) {
	var rowSlice []PaymentData
	rows, err := conn.Query(context.Background(),
		`SELECT row_to_json(row) FROM payment row`)
	if err != nil {
		return rowSlice, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		rows.Scan(&data)
		payment, err := JsonToPayment(data)
		if err != nil {
			return rowSlice, err
		}
		rowSlice = append(rowSlice, payment)
	}
	if err := rows.Err(); err != nil {
		return rowSlice, err
	}
	return rowSlice, nil
}

func GetAllOrders(conn *pgx.Conn) ([]OrderData, error) {
	var rowSlice []OrderData
	rows, err := conn.Query(context.Background(),
		`SELECT row_to_json(row) FROM orders row`)
	if err != nil {
		return rowSlice, err
	}
	defer rows.Close()

	for rows.Next() {
		var data []byte
		rows.Scan(&data)
		order, err := JsonToOrder(data)
		if err != nil {
			return rowSlice, err
		}
		rowSlice = append(rowSlice, order)
	}
	if err := rows.Err(); err != nil {
		return rowSlice, err
	}
	return rowSlice, nil
}

func InsertDelivery(tx pgx.Tx, delivery DeliveryData, OrderUID string) error {
	sqlInsertDelivery := `
		INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	_, err := tx.Exec(context.Background(),
		sqlInsertDelivery,
		OrderUID,
		delivery.Name, delivery.Phone, delivery.Zip, delivery.City,
		delivery.Address, delivery.Region, delivery.Email,
	)

	return err
}

func InsertPayment(tx pgx.Tx, payment PaymentData, OrderUID string) error {
	sqlInsertPayment := `
		INSERT INTO payment (order_uid, transaction, request_id, currency,
			provider, amount, payment_dt, bank, delivery_cost, goods_total, 
			custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	_, err := tx.Exec(context.Background(),
		sqlInsertPayment,
		OrderUID,
		payment.Transaction, payment.RequestId, payment.Currency,
		payment.Provider, payment.Amount, payment.PaymentDt, payment.Bank,
		payment.DeliveryCost, payment.GoodsTotal, payment.CustomFee,
	)

	return err
}

func InsertItem(tx pgx.Tx, item ItemData, OrderUID string) error {
	sqlInsertItem := `
		INSERT INTO items (order_uid, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	_, err := tx.Exec(context.Background(),
		sqlInsertItem,
		OrderUID,
		item.ChrtID, item.TrackNumber, item.Price, item.RID, item.Name,
		item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand,
		item.Status,
	)

	return err
}

func InsertOrder(order OrderData, conn *pgx.Conn) error {
	mutex.Lock()
	tx, err := conn.Begin(context.Background())

	if err != nil {
		return err
	}

	defer tx.Rollback(context.Background())

	sqlInsertOrder := `
		INSERT INTO orders (order_uid, track_number, entry, locale, 
			internal_signature, customer_id, delivery_service, shardkey, sm_id, 
			date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);`

	_, err = tx.Exec(context.Background(),
		sqlInsertOrder,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	err = InsertDelivery(tx, order.Delivery, order.OrderUID)
	if err != nil {
		return err
	}

	err = InsertPayment(tx, order.Payment, order.OrderUID)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		err = InsertItem(tx, item, order.OrderUID)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	mutex.Unlock()
	return nil
}
