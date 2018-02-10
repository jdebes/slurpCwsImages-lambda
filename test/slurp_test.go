package test

import (
	"github.com/jdebes/slurpCwsImages-lambda/model"
	"github.com/jdebes/slurpCwsImages-lambda/slurp"
	"github.com/jdebes/slurpCwsImages-lambda/test/mocks"
	. "github.com/onsi/ginkgo"
	"github.com/stretchr/testify/mock"
)

var _ = Describe("slurpCwsImages", func() {

	cwsService := new(mocks.CwsService)
	awsService := new(mocks.AwsService)

	var (
		productsResponse *model.ProductsResponse
	)

	BeforeEach(func() {
		image := model.Image{}
		image.Image = "https://api.codeswholesale.com/assets/images/no-image.jpg"
		image.Format = "medium"

		product := model.CwsProduct{}
		product.ProductID = "1234-ffff"
		product.Images = []model.Image{image}
		product.Regions = []string{"WORLDWIDE"}
		product.Name = "Test game"
		product.Platform = "Origin"

		productsResponse = &model.ProductsResponse{}
		productsResponse.Items = []model.CwsProduct{product}

		cwsService = new(mocks.CwsService)
		awsService = new(mocks.AwsService)

		cwsService.On("GetProducts").Return(productsResponse, nil)
		cwsService.On("GetProductImage", mock.Anything).Return([]byte{'t', 'e', 's', 't'}, ".jpg", nil)
		awsService.On("UploadItemToS3", []byte{'t', 'e', 's', 't'}, ".jpg", "1234-ffff", "medium").Return(nil)
	})

	Context("When calling SlurpImages with a valid image that does not exist in S3", func() {
		BeforeEach(func() {
			cwsService.On("HeadProductImageForUrl", mock.Anything).Return("https://api.codeswholesale.com/assets/images/test.jpg", nil)
			awsService.On("S3ItemExists", "1234-ffff", ".jpg", "medium").Return(false)

			slurp.SlurpImages(cwsService, awsService)
		})

		It("Should call UploadItemToS3", func() {
			awsService.AssertCalled(GinkgoT(), "UploadItemToS3", mock.Anything, ".jpg", "1234-ffff", "medium")
		})
	})

	Context("When calling SlurpImages with the no image url", func() {
		BeforeEach(func() {
			cwsService.On("HeadProductImageForUrl", mock.Anything).Return("https://api.codeswholesale.com/assets/images/no-image.jpg", nil)
			awsService.On("S3ItemExists", "1234-ffff", ".jpg", "medium").Return(false)

			slurp.SlurpImages(cwsService, awsService)
		})

		It("Should not call UploadItemToS3", func() {
			awsService.AssertNotCalled(GinkgoT(), "UploadItemToS3", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})

		It("Should not call S3ItemExists", func() {
			awsService.AssertNotCalled(GinkgoT(), "S3ItemExists", mock.Anything, mock.Anything, mock.Anything)
		})
	})

	Context("When calling SlurpImages with a valid image that exists in S3", func() {
		BeforeEach(func() {
			cwsService.On("HeadProductImageForUrl", mock.Anything).Return("https://api.codeswholesale.com/assets/images/test.jpg", nil)
			awsService.On("S3ItemExists", "1234-ffff", ".jpg", "medium").Return(true)

			slurp.SlurpImages(cwsService, awsService)
		})

		It("Should not call UploadItemToS3", func() {
			awsService.AssertNotCalled(GinkgoT(), "UploadItemToS3", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
		})
	})
})
