package filetype

import (
	"encoding/base64"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	HEADER_SIZE  = 512
	PNG_FILE     = "./test_data/pngImage.png"
	JPG_FILE     = "./test_data/jpgImage.jpg"
	JPEG_FILE    = "./test_data/jpegImage.jpeg"
	NO_EXT_PNG   = "./test_data/noExtensionPng"
	INVALID_FILE = "./test_data/wrongFile.txt"
)

func TestPNGfile(t *testing.T) {
	header, err := readImageHeaderToBase64(PNG_FILE)
	require.NoError(t, err)

	isValid, err := IsValid(header)

	require.NoError(t, err)
	require.True(t, isValid)
}

func TestJPGfile(t *testing.T) {
	header, err := readImageHeaderToBase64(JPG_FILE)
	require.NoError(t, err)

	isValid, err := IsValid(header)

	require.NoError(t, err)
	require.True(t, isValid)
}

func TestJPEGfile(t *testing.T) {
	header, err := readImageHeaderToBase64(JPEG_FILE)
	require.NoError(t, err)

	isValid, err := IsValid(header)

	require.NoError(t, err)
	require.True(t, isValid)
}

func TestValidImageWrongExtension(t *testing.T) {
	header, err := readImageHeaderToBase64(JPEG_FILE)
	require.NoError(t, err)

	isValid, err := IsValid(header)

	require.NoError(t, err)
	require.True(t, isValid)
}

func TestInvalidFile(t *testing.T) {
	header, err := readImageHeaderToBase64(INVALID_FILE)
	require.NoError(t, err)

	isValid, err := IsValid(header)

	require.NoError(t, err)
	require.False(t, isValid)
}

func readImageHeaderToBase64(imagePath string) (string, error) {
	header, err := readFileHeader(imagePath)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(header), nil
}

func readFileHeader(filepath string) ([]byte, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	header := make([]byte, HEADER_SIZE)
	_, err = r.Read(header)
	if err != nil {
		return nil, err
	}
	return header, nil
}
