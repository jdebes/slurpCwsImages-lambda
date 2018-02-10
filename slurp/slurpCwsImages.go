package slurp

import (
	"path/filepath"
	"slurpCwsImages/service"
	"strings"
)

const (
	noImageUrl = "https://api.codeswholesale.com/assets/images/no-image.jpg"
)

func SlurpImages(cwsService service.CwsService, awsService service.AwsService) {
	resp, err := cwsService.GetProducts()
	if err != nil {
		panic(1)
	}

	for _, item := range resp.Items {
		for _, region := range item.Regions {
			if strings.ToUpper(region) == "WORLDWIDE" {
				for _, image := range item.Images {
					fileUrl, err := cwsService.HeadProductImageForUrl(image.Image)
					if err != nil || fileUrl == noImageUrl {
						continue
					}
					fileExt := filepath.Ext(fileUrl)

					if !awsService.S3ItemExists(item.ProductID, fileExt, image.Format) {
						file, _, err := cwsService.GetProductImage(image.Image)
						if err != nil {
							continue
						}

						awsService.UploadItemToS3(file, fileExt, item.ProductID, image.Format)
					}
				}
			}
		}
	}
}
