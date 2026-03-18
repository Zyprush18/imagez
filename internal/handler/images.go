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
			return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed to convert image", err.Error(),  map[string]string{}))
		}
		if err.Error() == utils.UNSUPPORTED_FORMAT {
			return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed to convert image", fmt.Sprintf("%s: %s", err.Error(), formats[0]), map[string]string{}))
		}

		return c.JSON(http.StatusInternalServerError, utils.NewResponse("Failed to convert image", "Failed to convert image", map[string]string{}))
	}

	return c.JSON(http.StatusOK, utils.NewResponse("Image converted successfully", "",	 map[string]string{
		"file_name": fileName,
	}))
}

func (h *HandleImage) Resize(c *echo.Context) error {
	file, err := c.MultipartForm()
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed Request", "Failed to retrieve image file", map[string]string{}))
	}

	imgFile := file.File["images"]
	fileVal := file.Value

	if imgFile == nil || fileVal["width"] == nil || fileVal["height"] == nil {
		return c.JSON(http.StatusBadRequest,utils.NewResponse("Failed Request", "Missing required fields", map[string]string{}))
	}

	width, err := utils.ConvertToInt(fileVal["width"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed Request", "Invalid width value", map[string]string{}))
	}

	height, err := utils.ConvertToInt(fileVal["height"][0])
	if err != nil {
		return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed Request", "Invalid height value", map[string]string{}))
	}

	nameZipFile, err := h.svc.Resize(imgFile, width, height)

	if err != nil {
		c.Logger().Error(err.Error())

		if err.Error() == utils.UNSUPPORTED_TYPE {
			return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed to resize image", err.Error(),  map[string]string{}))
		}
		if err.Error() == utils.UNSUPPORTED_FORMAT {
			return c.JSON(http.StatusBadRequest, utils.NewResponse("Failed to resize image", "image file format not supported", map[string]string{}))
		}
		return c.JSON(http.StatusInternalServerError, utils.NewResponse("Failed to resize image", err.Error(), map[string]string{}))
	}

	return c.JSON(http.StatusOK, utils.NewResponse("Image resized successfully", "", map[string]string{
		"file_name": nameZipFile,
	}))
}

func (h *HandleImage) Compress(c *echo.Context) error {
	return c.JSON(http.StatusOK, utils.NewResponse("Image compressed successfully", "", map[string]string{}))
}
