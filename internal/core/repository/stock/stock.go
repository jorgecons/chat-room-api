package stock

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-resty/resty/v2"
)

type (
	Repo struct {
		restyClient *resty.Client
	}
)

func NewStockPriceRepo(restyClient *resty.Client) *Repo {
	return &Repo{
		restyClient: restyClient,
	}
}

func (r *Repo) GetPrice(ctx context.Context, stockID string) (float64, error) {
	resp, err := r.restyClient.R().SetContext(ctx).Get(fmt.Sprintf(r.restyClient.BaseURL, stockID))
	if err != nil {
		return 0, err
	}
	if resp.StatusCode() >= http.StatusBadRequest {
		return 0, errors.New("error")
	}
	fmt.Println(string(resp.Body()))
	reader := csv.NewReader(strings.NewReader(resp.String()))
	_, err = reader.Read()
	if err != nil {
		log.Println("CSV Read Error:", err)
		return 0, errors.New("invalid stock code")
	}

	record, err := reader.Read()
	if err != nil {
		return 0, errors.New("invalid stock code")
	}

	stockPrice := record[6]
	stockPriceFloat, _ := strconv.ParseFloat(stockPrice, 64)

	return stockPriceFloat, nil
}
