package server

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

type Server struct {
	e              *echo.Echo
	imageInterface ImageInterface
}

type ImageInterface interface {
	DrawText(inputFileName, topText, bottomText string) (string, error)
}

func NewServer(imageInterface ImageInterface) *Server {
	e := echo.New()
	e.POST("/upload", func(c echo.Context) error {
		// Bind the JSON metadata to a struct
		metadata := new(Metadata)

		topText := c.FormValue("top_text")
		bottomText := c.FormValue("bottom_text")
		metadata.TopText = topText
		metadata.BottomText = bottomText

		// Get the file from the request
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "File upload failed"})
		}

		// Open the file for reading
		src, err := file.Open()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to open the file"})
		}
		defer src.Close()

		dst, err := os.CreateTemp("", "sample")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create temporary file"})
		}
		defer dst.Close()

		if _, err = dst.ReadFrom(src); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to save the file"})
		}

		resultFileName, err := imageInterface.DrawText(dst.Name(), topText, bottomText)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to create temporary file"})
		}
		// Respond with success and metadata
		return c.File(resultFileName)
	})

	return &Server{e: e}
}

func (s *Server) Start() {
	s.e.Start(":8080")
}

func (s *Server) Stop() {
	s.e.Shutdown(context.Background())
}

type Metadata struct {
	TopText    string `json:"top_text"`
	BottomText string `json:"bottom_text"`
}
