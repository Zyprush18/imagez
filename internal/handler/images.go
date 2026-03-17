package handler

import (
	"fmt"
	"net/http"

	"github.com/Zyprush18/imagez/internal/service"
	"github.com/Zyprush18/imagez/utils"
	"github.com/labstack/echo/v5"
)

type HandleImage struct {
	svc service.ImagesService
}

func NewHandleImage(svc service.ImagesService) *HandleImage {
	return &HandleImage{svc: svc}
}

func (h *HandleImage) Convert(c *echo.Context) error {
	file, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Failed to retrieve image file",
		})
	}

	files := file.File["images"]
	formats := file.Value["format"]

	if files == nil || formats == nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing field image or format",
		})
	}

	fileName, err := h.svc.Convert(files, formats[0])

	if err != nil {
		c.Logger().Error(err.Error())

		if err.Error() == utils.UNSUPPORTED_TYPE {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("%s: %s", err.Error(), formats[0]),
			})
		}
		if err.Error() == utils.UNSUPPORTED_FORMAT {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("%s: %s", err.Error(), formats[0]),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to convert image",
		})
	}

	return c.JSON(http.StatusOK, utils.NewResponse("Image converted successfully", map[string]string{
		"file_name": fileName,
	}))
}

func (h *HandleImage) Resize(c *echo.Context) error {
	return c.JSON(http.StatusOK, utils.NewResponse("Image resized successfully", "ok"))
}

func (h *HandleImage) Compress(c *echo.Context) error {
	return c.JSON(http.StatusOK, utils.NewResponse("Image compressed successfully", "ok"))
}
