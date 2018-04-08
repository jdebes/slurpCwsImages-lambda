package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"

	"github.com/jdebes/slurpCwsImages-lambda/model"
	log "github.com/sirupsen/logrus"
)

const (
	productsAPIPath = "https://api.codeswholesale.com/v1/products"
)

type CwsService interface {
	GetProducts() (*model.ProductsResponse, error)
	GetProductImage(path string) ([]byte, string, error)
	HeadProductImageForUrl(path string) (string, error)
}

type cwsServiceImpl struct {
	client           *http.Client
	noRedirectClient *http.Client
}

func (service *cwsServiceImpl) GetProducts() (*model.ProductsResponse, error) {
	var products model.ProductsResponse

	resp, err := service.client.Get(productsAPIPath)
	if err != nil || checkStatus(resp) != nil {
		logField := log.WithError(err)

		if err = checkStatus(resp); err != nil {
			logField = log.WithError(err)
			return nil, err
		}

		logField.Error("Failed to Get Products")
		return nil, err
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
		logField := log.WithError(err)

		if err = checkStatus(resp); err != nil {
			logField = log.WithError(err)
			return nil, "", err
		}

		logField.Error(fmt.Sprintf("Failed to Get %s", path))
		return nil, "", err
	}
	defer resp.Body.Close()

	var file []byte
	finalURL := resp.Request.URL.String()
	if filepath.Ext(finalURL) == ".png" {
		file, err = convertPngToJpg(resp.Body)
	} else {
		file, err = ioutil.ReadAll(resp.Body)
	}

	return file, FileExtension, nil
}

func (service *cwsServiceImpl) HeadProductImageForUrl(path string) (string, error) {
	var finalURL string
	resp, err := service.noRedirectClient.Get(path)
	if err != nil {
		if resp != nil && resp.StatusCode == 302 {
			finalURL = resp.Header.Get("Location")
			return finalURL, nil
		}

		log.WithError(err).Error(fmt.Sprintf("Failed to Head %s", path))
		if resp != nil {
			return "", checkStatus(resp)
		}

		return "", err
	}
	defer resp.Body.Close()

	finalURL = resp.Request.URL.String()
	log.WithField("finalUrl", finalURL).Info("Didn't get 302 redirect.")
	return finalURL, nil
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

func convertPngToJpg(file io.Reader) ([]byte, error) {
	image, err := png.Decode(file)
	if err != nil {
		return nil, err
	}

	buffer := new(bytes.Buffer)
	err = jpeg.Encode(buffer, image, nil)
	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
