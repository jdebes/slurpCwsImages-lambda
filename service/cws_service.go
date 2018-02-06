package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"errors"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"slurpCwsImages/model"
)

const (
	productsAPIPath = "https://api.codeswholesale.com/v1/products"
)

type CwsService interface {
	GetProducts() (*model.ProductsResponse, error)
	GetProductImage(path string) ([]byte, string, error)
	HeadProductImage(path string) (string, error)
}

type cwsServiceImpl struct {
	client           *http.Client
	noRedirectClient *http.Client
}

func (service *cwsServiceImpl) GetProducts() (*model.ProductsResponse, error) {
	var products model.ProductsResponse

	resp, err := service.client.Get(productsAPIPath)
	if err != nil || checkStatus(resp) != nil {
		log.WithError(err).Error("Failed to Get Products")
		return nil, checkStatus(resp)
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&products)
	if err != nil {
		log.WithError(err).Error("Failed to unmarshal products response")
		return nil, err
	}

	return &products, nil
}

func (service *cwsServiceImpl) GetProductImage(path string) ([]byte, string, error) {
	resp, err := http.Get(path)
	if err != nil || checkStatus(resp) != nil {
		log.WithError(err).Error(fmt.Sprintf("Failed to Get %s", path))
		return nil, "", checkStatus(resp)
	}
	defer resp.Body.Close()

	file, err := ioutil.ReadAll(resp.Body)

	finalURL := resp.Request.URL.String()

	return file, filepath.Ext(finalURL), nil
}

func (service *cwsServiceImpl) HeadProductImage(path string) (string, error) {
	var finalURL string
	resp, err := service.noRedirectClient.Get(path)
	if err != nil {
		if resp.StatusCode == 302 {
			finalURL = resp.Header.Get("Location")
			return filepath.Ext(finalURL), nil
		}

		log.WithError(err).Error(fmt.Sprintf("Failed to Head %s", path))
		return "", checkStatus(resp)
	}
	defer resp.Body.Close()

	log.Error("Didn't get 302 redirect. Check CWS API")

	finalURL = resp.Request.URL.String()
	return filepath.Ext(finalURL), nil
}

func BuildCwsService(cwsClient *http.Client) CwsService {
	return &cwsServiceImpl{
		client: cwsClient,
		noRedirectClient: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return errors.New("redirected")
			},
		},
	}
}
