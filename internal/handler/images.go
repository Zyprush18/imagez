package handler

import (
	"fmt"
	"net/http"

	"github.com/Zyprush18/imagez/internal/service"
	"github.com/Zyprush18/imagez/pkg"
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

	if err := h.svc.Convert(files, formats[0]);err != nil {
		if err.Error() == pkg.UNSUPPORTED_FORMAT {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("%s: %s", err.Error(), formats[0]),
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to convert image",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "ini api untuk image convert",
	})
}

func (h *HandleImage) Resize(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ini api untuk image resize",
	})
}

func (h *HandleImage) Compress(c *echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{
		"message": "ini api untuk image resize",
	})
}