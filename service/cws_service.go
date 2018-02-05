package service

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"slurpCwsImages/model"
	"path/filepath"
)

const (
	productsAPIPath = "https://api.codeswholesale.com/v1/products"
)

type CwsService interface {
	GetProducts() (*model.ProductsResponse, error)
	GetProductImage(path string) ([]byte, string, error)
}

type cwsServiceImpl struct {
	client *http.Client
}

func (service *cwsServiceImpl) GetProducts() (*model.ProductsResponse, error) {
	var products model.ProductsResponse

	resp, err := service.client.Get(productsAPIPath)
	if err != nil || checkStatus(resp) != nil {
		log.WithError(err).Error("Failed to Get Products")
		return nil, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal products response")
	}

	return &products, nil
}

func (service *cwsServiceImpl) GetProductImage(path string) ([]byte, string, error) {
	resp, err := http.Get(path)
	if err != nil || checkStatus(resp) != nil {
		log.WithError(err).Error(fmt.Sprintf("Failed to Get %s", path))
		return nil, "", err
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)

	finalURL := resp.Request.URL.String()

	return file, filepath.Ext(finalURL), nil
}

func BuildCwsService(cwsClient *http.Client) CwsService {
	return &cwsServiceImpl{cwsClient}
}




