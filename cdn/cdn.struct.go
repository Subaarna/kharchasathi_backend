package cdn

import (
	"github.com/cloudinary/cloudinary-go/api/admin"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/cloudinary/cloudinary-go/config"
	"github.com/cloudinary/cloudinary-go/logger"
)

type Cloudinary struct {
	Config config.Configuration
	Admin  admin.API
	Upload uploader.API
	Logger *logger.Logger
}
