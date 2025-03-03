package stock

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"chat-room-api/internal/core/domain"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"
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
	reader := csv.NewReader(strings.NewReader(resp.String()))
	_, err = reader.Read()
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("stock", stockID).Error("error getting stock price")
		return 0, domain.WrapError(domain.ErrGettingStockPrice, err)
	}

	record, err := reader.Read()
	if err != nil {
		logrus.WithContext(ctx).WithError(err).WithField("stock", stockID).Error("error reading stock price")
		return 0, domain.WrapError(domain.ErrGettingStockPrice, err)
	}

	stockPrice := record[6]
	if stockPrice == nilValue {
		logrus.WithContext(ctx).WithField("stock", stockID).Error("stock price is nil")
		return 0, domain.ErrStockPriceNotFound
	}
	stockPriceFloat, _ := strconv.ParseFloat(stockPrice, 64)

	return stockPriceFloat, nil
}
