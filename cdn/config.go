package cdn

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/logger"
)

func NewFromConfiguration(configuration config.Configuration) (*Cloudinary, error) {
	logger := logger.New()

	return &Cloudinary{
		Config: configuration,
		Admin: admin.API{
			Config: configuration,
			Logger: logger,
		},
		Upload: uploader.API{
			Config: configuration,
			Logger: logger,
		},
		Logger: logger,
	}, nil
}

func CdnSetting() (*Cloudinary, error) {

	cld, err := cloudinary.NewFromParams(os.Getenv("CLOUD_NAME"), os.Getenv("CLOUD_KEY"), os.Getenv("CLOUD_API_SECRET"))
	if err != nil {
		return nil, err
	}
	// fmt.Println("cld: ", cld)
	return NewFromConfiguration(cld.Config)
}
