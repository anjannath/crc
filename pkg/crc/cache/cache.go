package cache

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/code-ready/crc/pkg/crc/constants"
	"github.com/code-ready/crc/pkg/crc/logging"
	"github.com/code-ready/crc/pkg/download"
	"github.com/code-ready/crc/pkg/embed"
	"github.com/code-ready/crc/pkg/extract"
	crcos "github.com/code-ready/crc/pkg/os"
	"github.com/pkg/errors"
)

type BinaryCache struct {
	binaryName string
	archiveURL string
	destDir    string
}

func New(binaryName string, archiveURL string, destDir string) *BinaryCache {
	return &BinaryCache{binaryName: binaryName, archiveURL: archiveURL, destDir: destDir}
}

func NewOcCache(destDir string) *BinaryCache {
	return New(constants.OcBinaryName, constants.GetOcUrl(), destDir)
}

func NewPodmanCache(destDir string) *BinaryCache {
	return New(constants.PodmanBinaryName, constants.GetPodmanUrl(), destDir)
}

func NewGoodhostsCache(destDir string) *BinaryCache {
	return New(constants.GoodhostsBinaryName, constants.GetGoodhostsUrl(), destDir)
}

func (c *BinaryCache) IsCached() bool {
	if _, err := os.Stat(filepath.Join(c.destDir, c.binaryName)); os.IsNotExist(err) {
		return false
	}
	return true
}

func (c *BinaryCache) EnsureIsCached() error {
	if !c.IsCached() {
		err := c.DoCache()
		if err != nil {
			return err
		}
	}
	return nil
}

// CacheBinary downloads and caches the requested binary into the CRC directory
func (c *BinaryCache) DoCache() error {
	if c.IsCached() {
		return nil
	}

	// Create tmp dir to download the requested tarball
	tmpDir, err := ioutil.TempDir("", "crc")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	assetTmpFile, err := c.extract(tmpDir)
	if err != nil {
		return err
	}

	// Extract the tarball and put it the cache directory.
	extractedFiles, err := extract.UncompressWithFilter(assetTmpFile, tmpDir,
		func(filename string) bool { return filepath.Base(filename) == c.binaryName })
	if err != nil {
		return errors.Wrapf(err, "Cannot uncompress '%s'", assetTmpFile)
	}

	// Copy the requested asset into its final destination
	err = os.MkdirAll(c.destDir, 0750)
	if err != nil && !os.IsExist(err) {
		return errors.Wrap(err, "Cannot create the target directory.")
	}

	for _, extractedFilePath := range extractedFiles {
		finalBinaryPath := filepath.Join(c.destDir, filepath.Base(extractedFilePath))
		err = crcos.CopyFileContents(extractedFilePath, finalBinaryPath, 0500)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *BinaryCache) extract(destDir string) (string, error) {
	logging.Debug("Trying to extract oc from crc binary")
	archiveName := filepath.Base(c.archiveURL)
	destPath := filepath.Join(destDir, archiveName)
	err := embed.Extract(archiveName, destPath)
	if err != nil {
		logging.Debug("Downloading oc")
		return download.Download(c.archiveURL, destDir, 0600)
	}

	return destPath, err
}
