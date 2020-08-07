// Package core represents the core functionality of all commands
package core

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gitrob/common"
)

// Set various internal values used by the web interface
const (
	GithubBaseUri   = "https://raw.githubusercontent.com"
	MaximumFileSize = 153600
	GitLabBaseUri   = "https://gitlab.com"
	CspPolicy       = "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'"
	ReferrerPolicy  = "no-referrer"
)

// Is this a github repo/org
var IsGithub bool

// binaryFileSystem  holds a filesystem handle
type binaryFileSystem struct {
	fs http.FileSystem
}

// Open will return an http file object that refers to a given file
func (b *binaryFileSystem) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

// Exists checks if a given file with a given prefix exists and attempts to open it
func (b *binaryFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

// BinaryFileSystem returns a binary file system object used by the web frontend
func BinaryFileSystem(root string) *binaryFileSystem {
	fs := &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: AssetInfo, Prefix: root}
	return &binaryFileSystem{
		fs,
	}
}

// NewRouter will create an instance of the web frontend, setting the necessary parameters.
func NewRouter(s *Session) *gin.Engine {

	if s.ScanType == "github" {
		IsGithub = true
	}

	if s.Debug == true {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.Use(static.Serve("/", BinaryFileSystem("static")))
	router.Use(secure.New(secure.Config{
		SSLRedirect:           false,
		IsDevelopment:         false,
		FrameDeny:             true,
		ContentTypeNosniff:    true,
		BrowserXssFilter:      true,
		ContentSecurityPolicy: CspPolicy,
		ReferrerPolicy:        ReferrerPolicy,
	}))
	router.GET("/stats", func(c *gin.Context) {
		c.JSON(200, s.Stats)
	})
	router.GET("/findings", func(c *gin.Context) {
		c.JSON(200, s.Findings)
	})
	router.GET("/targets", func(c *gin.Context) {
		c.JSON(200, s.Targets)
	})
	router.GET("/repositories", func(c *gin.Context) {
		c.JSON(200, s.Repositories)
	})
	router.GET("/files/:owner/:repo/:commit/*path", fetchFile)

	return router
}

// TODO this will fail for other target types and must be converted to a switch for scalability
// fetchFile returns a given path to a file that can be cicked on by a user
func fetchFile(c *gin.Context) {
	fileUrl := func() string {
		if IsGithub {
			return fmt.Sprintf("%s/%s/%s/%s%s", GithubBaseUri, c.Param("owner"), c.Param("repo"), c.Param("commit"), c.Param("path"))
		} else {
			results := common.CleanUrlSpaces(c.Param("owner"), c.Param("repo"), c.Param("commit"), c.Param("path"))
			return fmt.Sprintf("%s/%s/%s/%s/%s%s", GitLabBaseUri, results[0], results[1], "/-/raw/", results[2], results[3])
		}
	}()
	resp, err := http.Head(fileUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	if resp.StatusCode == http.StatusNotFound {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No content",
		})
		return
	}

	if resp.ContentLength > MaximumFileSize {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": fmt.Sprintf("File size exceeds maximum of %d bytes", MaximumFileSize),
		})
		return
	}

	resp, err = http.Get(fileUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	c.String(http.StatusOK, string(body[:]))
}
