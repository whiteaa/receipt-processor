package db

import (
	"errors"
	"sync"
)

type db struct {
	m sync.Map
}

var once sync.Once
var ReceiptDB db

var ErrReceiptNotFound = errors.New("receipt not found")
var ErrInvalidPointsType = errors.New("invalid points type")

func InitDB() {
	once.Do(createDB)
}

func createDB() {
	ReceiptDB = db{
		m: sync.Map{},
	}
}

func (r *db) GetReceipt(receiptId string) (int, error) {
	val, ok := r.m.Load(receiptId)
	if !ok {
		return 0, ErrReceiptNotFound
	}
	points, ok := val.(int)
	if !ok {
		return 0, ErrInvalidPointsType
	}
	return points, nil
}

func (r *db) SetReceipt(receiptId string, points int) {
	r.m.Store(receiptId, points)
}
