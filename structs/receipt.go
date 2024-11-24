package structs

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type ProcessReceiptsV1Response struct {
	Id string `json:"id"`
}

type GetPointsV1Response struct {
	Points int `json:"points"`
}

type Receipt struct {
	Retailer     string  `json:"retailer" validate:"required"`
	PurchaseDate string  `json:"purchaseDate" validate:"required"`
	PurchaseTime string  `json:"purchaseTime" validate:"required"`
	Items        []Items `json:"items" validate:"required"`
	Total        string  `json:"total" validate:"required"`
}

type Items struct {
	ShortDescription string `json:"shortDescription" validate:"required"`
	Price            string `json:"price" validate:"required"`
}

func (r *Receipt) CalculatePoints() int {
	points := r.retailerPoints()
	points += r.totalPoints()
	points += r.itemsPoints()
	points += r.purchaseDatePoints()
	points += r.purchaseTimePoints()
	return points
}

func (r *Receipt) retailerPoints() (points int) {
	// One point for every alphanumeric character in the retailer name
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Errorf("Error compiling regex err: %v", err)
		return points
	}
	result := reg.ReplaceAllString(r.Retailer, "")
	points += len(result)

	return points
}

func (r *Receipt) totalPoints() (points int) {
	total, err := strconv.ParseFloat(r.Total, 64)
	if err != nil {
		log.Errorf("Error converting total: %s err: %v", r.Total, err)
		return 0
	}

	// 50 points if the total is a round dollar amount with no cents
	if math.Trunc(total) == total {
		points += 50
	}
	// 25 points if the total is a multiple of 0.25
	if math.Mod(total, 0.25) == 0 {
		points += 25
	}
	return points
}

func (r *Receipt) itemsPoints() (points int) {
	// 5 points for every two items on the receipt
	points += 5 * (len(r.Items) / 2)

	// If the trimmed length of the item description is a multiple of 3,
	// multiply the price by 0.2 and round up to the nearest integer.
	// The result is the number of points earned.
	for _, item := range r.Items {
		trimmed := strings.TrimSpace(item.ShortDescription)
		if len(trimmed)%3 == 0 {
			price, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				log.Errorf("Error converting price: %s err: %v", item.Price, err)
				continue
			}
			points += int(math.Ceil(0.2 * price))
		}
	}

	return points
}

func (r *Receipt) purchaseDatePoints() (points int) {
	t, err := time.Parse("2006-01-02", r.PurchaseDate)
	if err != nil {
		log.Errorf("Error converting purchaseDate: %s err: %v", r.PurchaseDate, err)
		return points
	}

	// 6 points if the day in the purchase date is odd.
	if t.Day()%2 == 1 {
		points += 6
	}

	return points
}

func (r *Receipt) purchaseTimePoints() (points int) {
	t, err := time.Parse("15:04", r.PurchaseTime)
	if err != nil {
		log.Errorf("Error converting purchaseTime: %s err: %v", r.PurchaseTime, err)
		return points
	}

	// 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	twoPM := time.Date(0, 1, 1, 14, 0, 0, 0, time.UTC)
	fourPM := time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC)
	if t.After(twoPM) && t.Before(fourPM) {
		points += 10
	}
	return points
}
