package handler

import (
	"net/http"
	"receiptprocessor/db"
	"receiptprocessor/structs"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

// GetPointsV1 looks up the receipt by the ID and returns an object specifying
// the points awarded
func GetPointsV1(c echo.Context) error {
	receiptId := c.Param("id")

	points, err := db.ReceiptDB.GetReceipt(receiptId)
	if err != nil {
		log.Errorf("Invalid receiptId: %s err: %v", receiptId, err.Error())
		return echo.NewHTTPError(http.StatusNotFound)
	}
	resp := structs.GetPointsV1Response{
		Points: points,
	}
	return c.JSON(http.StatusOK, &resp)
}

// ProcessReceiptsV1 Takes in a JSON receipt and returns a JSON object with
// a generated ID. The points are calculated and stored by the ID
func ProcessReceiptsV1(c echo.Context) error {
	var receipt structs.Receipt
	if err := (&echo.DefaultBinder{}).BindBody(c, &receipt); err != nil {
		log.Errorf("Error binding request err: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid body")
	}
	if err := c.Validate(receipt); err != nil {
		log.Errorf("Invalid request err: %v", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	receiptId := uuid.NewString()
	points := receipt.CalculatePoints()

	db.ReceiptDB.SetReceipt(receiptId, points)

	resp := structs.ProcessReceiptsV1Response{
		Id: receiptId,
	}
	return c.JSON(http.StatusOK, &resp)
}
