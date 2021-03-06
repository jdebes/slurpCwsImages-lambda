package slurp

import (
	"strings"
	"sync"

	"github.com/jdebes/slurpCwsImages-lambda/model"
	"github.com/jdebes/slurpCwsImages-lambda/service"
	"github.com/sirupsen/logrus"
)

const (
	noImageUrl = "https://api.codeswholesale.com/assets/images/no-image.jpg"
)

func SlurpImages(cwsService service.CwsService, awsService service.AwsService, concurrencyLimit int) {
	resp, err := cwsService.GetProducts()
	if err != nil {
		panic(1)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrencyLimit)
	uploadCount := 0

	for _, item := range resp.Items {
		for _, region := range item.Regions {
			if strings.ToUpper(region) == "WORLDWIDE" {
				for _, image := range item.Images {
					wg.Add(1)
					sem <- struct{}{}
					go func(image model.Image, item model.CwsProduct) {
						uploadImage(image, item, cwsService, awsService)
						uploadCount++
						defer wg.Done()
						defer func() { <-sem }()
					}(image, item)
				}
			}
		}
	}

	wg.Wait()
	logrus.WithField("count", uploadCount).Info("Slurp completed")
}

func uploadImage(image model.Image, item model.CwsProduct, cwsService service.CwsService, awsService service.AwsService) {
	fileUrl, err := cwsService.HeadProductImageForUrl(image.Image)
	if err != nil || fileUrl == noImageUrl {
		logrus.WithField("cwsId", item.ProductID).Info("Did not upload product image")
		return
	}

	if !awsService.S3ItemExists(item.ProductID, image.Format) {
		file, fileExt, err := cwsService.GetProductImage(image.Image)
		if err != nil {
			return
		}

		awsService.UploadItemToS3(file, fileExt, item.ProductID, image.Format)
	}
}
