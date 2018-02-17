package slurp

import (
	"path/filepath"
	"strings"
	"sync"

	"github.com/jdebes/slurpCwsImages-lambda/model"
	"github.com/jdebes/slurpCwsImages-lambda/service"
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
	wg.Add(len(resp.Items))

	sem := make(chan struct{}, concurrencyLimit)

	for _, item := range resp.Items {
		for _, region := range item.Regions {
			if strings.ToUpper(region) == "WORLDWIDE" {
				for _, image := range item.Images {
					sem <- struct{}{}
					go func() {
						uploadImage(image, item, cwsService, awsService)
						defer wg.Done()
						defer func() { <-sem }()
					}()
				}
			}
		}
	}

	wg.Wait()
}

func uploadImage(image model.Image, item model.CwsProduct, cwsService service.CwsService, awsService service.AwsService) {
	fileUrl, err := cwsService.HeadProductImageForUrl(image.Image)
	if err != nil || fileUrl == noImageUrl {
		return
	}
	fileExt := filepath.Ext(fileUrl)

	if !awsService.S3ItemExists(item.ProductID, fileExt, image.Format) {
		file, _, err := cwsService.GetProductImage(image.Image)
		if err != nil {
			return
		}

		awsService.UploadItemToS3(file, fileExt, item.ProductID, image.Format)
	}
}
