package infrastructure

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/Nickolasll/gomart/internal/domain"
)

type AccrualClient struct {
	URL string
}

func (c AccrualClient) GetOrderStatus(number string) (domain.AccrualOrderResponse, error) {
	var accrualResponse domain.AccrualOrderResponse
	client := &http.Client{}
	req, _ := http.NewRequest("GET", c.URL+"/"+number, nil)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil || res.StatusCode > http.StatusNoContent {
		return accrualResponse, domain.ErrAccrualIsBusy
	}
	if res.StatusCode == http.StatusNoContent {
		return accrualResponse, domain.ErrDocumentNotFound
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return accrualResponse, err
	}
	err = json.Unmarshal(body, &accrualResponse)
	if err != nil {
		return accrualResponse, err
	}
	return accrualResponse, nil
}
