package imagestore

import (
	"bytes"
	"errors"
	"image"
	"os"
	"path/filepath"
	"strings"

	"github.com/deepteams/webp"
	"github.com/disintegration/imaging"
	"github.com/rs/zerolog"
)

var ErrUnsupportedFormat = errors.New("unsupported or corrupted image file")

type SavedImage struct {
	FileName string
	ImageURL string
	Width    int
	Height   int
}

type Store interface {
	Save(referenceType, baseName string, data []byte) (*SavedImage, error)
	Remove(referenceType string, fileNames ...string) error
}

type LocalStore struct {
	baseDir  string
	baseURL  string
	maxWidth int
	logger   zerolog.Logger
}

func NewLocalStore(baseDir, baseURL string, maxWidth int, logger zerolog.Logger) *LocalStore {
	return &LocalStore{
		baseDir:  baseDir,
		baseURL:  strings.TrimRight(baseURL, "/"),
		maxWidth: maxWidth,
		logger:   logger,
	}
}

// Save resizes the upload and writes only a single .webp file to disk.
// The original upload is never persisted. baseName is the desired file name
// without extension (e.g. "BP1-FR-00125-a8k2x").
func (s *LocalStore) Save(referenceType, baseName string, data []byte) (*SavedImage, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, ErrUnsupportedFormat
	}

	if s.maxWidth > 0 && img.Bounds().Dx() > s.maxWidth {
		img = imaging.Resize(img, s.maxWidth, 0, imaging.Lanczos)
	}

	dir := filepath.Join(s.baseDir, referenceType)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	webpName := filepath.Base(strings.TrimSpace(baseName)) + ".webp"
	webpPath := filepath.Join(dir, webpName)

	webpFile, err := os.Create(webpPath)
	if err != nil {
		return nil, err
	}

	if err := webp.Encode(webpFile, img, nil); err != nil {
		webpFile.Close()
		s.removeFiles(webpPath)
		return nil, err
	}

	if err := webpFile.Close(); err != nil {
		s.removeFiles(webpPath)
		return nil, err
	}

	bounds := img.Bounds()
	return &SavedImage{
		FileName: webpName,
		ImageURL: s.baseURL + "/" + referenceType + "/" + webpName,
		Width:    bounds.Dx(),
		Height:   bounds.Dy(),
	}, nil
}

func (s *LocalStore) Remove(referenceType string, fileNames ...string) error {
	var firstErr error
	for _, name := range fileNames {
		if name == "" {
			continue
		}
		path := filepath.Join(s.baseDir, referenceType, filepath.Base(name))
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			s.logger.Error().Err(err).Str("path", path).Msg("failed to remove image file")
			if firstErr == nil {
				firstErr = err
			}
		}
	}
	return firstErr
}

func (s *LocalStore) removeFiles(paths ...string) {
	for _, path := range paths {
		if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
			s.logger.Warn().Err(err).Str("path", path).Msg("failed to clean up image file")
		}
	}
}
