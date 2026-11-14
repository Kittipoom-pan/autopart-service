package validation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Kittipoom-pan/autopart-service/internal/module/image/entity"
	customerror "github.com/Kittipoom-pan/autopart-service/pkg/error"
)

var allowedExtensions = map[string]struct{}{
	".jpg":  {},
	".jpeg": {},
	".png":  {},
	".webp": {},
}

func ValidateUploadImageRequest(req *entity.UploadImageReq, fileName string, sizeBytes, maxSizeBytes int64) error {
	errs := map[string]string{}

	if !entity.ReferenceType(req.ReferenceType).IsValid() {
		errs["reference_type"] = "unsupported reference type"
	}

	if req.ReferenceID <= 0 {
		errs["reference_id"] = "reference_id is required and must be greater than zero"
	}

	altText := strings.TrimSpace(req.AltText)
	if altText == "" {
		errs["alt_text"] = "alt_text is required"
	} else if len([]rune(altText)) > 150 {
		errs["alt_text"] = "alt_text must not exceed 150 characters"
	}

	switch {
	case strings.TrimSpace(fileName) == "":
		errs["file"] = "file is required"
	default:
		if _, ok := allowedExtensions[strings.ToLower(filepath.Ext(fileName))]; !ok {
			errs["file"] = "file must be one of .jpg, .jpeg, .png, .webp"
		} else if sizeBytes <= 0 {
			errs["file"] = "file is empty"
		} else if sizeBytes > maxSizeBytes {
			errs["file"] = fmt.Sprintf("file size must not exceed %dMB", maxSizeBytes/(1024*1024))
		}
	}

	if len(errs) > 0 {
		return customerror.InvalidRequestData(errs)
	}

	return nil
}
