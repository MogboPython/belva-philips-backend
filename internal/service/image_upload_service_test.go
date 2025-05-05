package service

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// helper to create a mock multipart.FileHeader
func createMultipartFileHeader(t *testing.T, filePath, contentType string) *multipart.FileHeader {
	file, err := os.Open(filePath)
	assert.NoError(t, err)
	defer file.Close()

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	assert.NoError(t, err)

	_, err = io.Copy(part, file)
	assert.NoError(t, err)

	err = writer.Close()
	assert.NoError(t, err)

	req := bytes.NewReader(b.Bytes())
	reader := multipart.NewReader(req, writer.Boundary())
	form, err := reader.ReadForm(10 << 20) // 10 MB max memory
	assert.NoError(t, err)

	fileHeaders := form.File["file"]
	assert.NotEmpty(t, fileHeaders)

	// manually override content-type if provided
	if contentType != "" {
		fileHeaders[0].Header.Set("Content-Type", contentType)
	}

	return fileHeaders[0]
}

func TestUploadAndRemoveImage(t *testing.T) {
	var publicURL string

	var err error

	t.Run("Should save image in supabase bucket", func(t *testing.T) {
		t.Parallel()

		filePath := "./test.png"
		contentType := "image/png"
		bucketID := "blog-cover-photos"

		imageFile := createMultipartFileHeader(t, filePath, contentType)
		publicURL, err = uploadFile(imageFile, bucketID)

		assert.NoError(t, err)
		assert.NotEmpty(t, publicURL)
		assert.True(t, strings.Contains(publicURL, "blog-cover-photos"))
	})

	// t.Run("Should remove image from supabase bucket", func(t *testing.T) {
	// 	if publicURL == "" {
	// 		t.Fatal("publicURL is empty - upload test probably failed")
	// 	}

	// 	err := removeFile(publicURL)

	// 	assert.NoError(t, err)
	// })
}
