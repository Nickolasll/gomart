package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nickolasll/gomart/internal/domain"
	"github.com/sirupsen/logrus"
)

type AccrualClient struct {
	URL string
	Log *logrus.Logger
}

func (c AccrualClient) GetOrderStatus(number string) (domain.AccrualOrderResponse, error) {
	c.Log.Info("requesting for: " + number)
	var accrualResponse domain.AccrualOrderResponse
	client := &http.Client{}
	req, err := http.NewRequest("GET", c.URL+"/"+number, nil)
	if err != nil {
		c.Log.Info("new request error " + err.Error())
		return accrualResponse, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil || res.StatusCode > http.StatusNoContent {
		c.Log.Info("res.StatusCode > http.StatusNoContent " + err.Error())
		return accrualResponse, domain.ErrAccrualIsBusy
	}
	if res.StatusCode == http.StatusNoContent {
		c.Log.Info("http.StatusNoContent " + err.Error())
		return accrualResponse, domain.ErrDocumentNotFound
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.Log.Info("io.ReadAll " + err.Error())
		return accrualResponse, err
	}
	err = json.Unmarshal(body, &accrualResponse)
	if err != nil {
		c.Log.Info("Unmarshal " + err.Error())
		return accrualResponse, err
	}
	return accrualResponse, nil
}
