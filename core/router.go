package core

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-contrib/secure"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
)

// Set various internal values used by the web interface
const (
	GithubBaseURI   = "https://raw.githubusercontent.com"
	MaximumFileSize = 153600
	GitLabBaseURL   = "https://gitlab.com"
	CspPolicy       = "default-src 'none'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'"
	ReferrerPolicy  = "no-referrer"
)

// Is this a github repo/org
var isGithub bool

// NewRouter will create an instance of the web frontend, setting the necessary parameters.
func NewRouter(s *Session) *gin.Engine {

	if s.ScanType == "github" {
		isGithub = true
	}

	if s.Debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	embedFS, err := static.EmbedFolder(staticFS, "static")
	if err != nil {
		s.Out.Fatal("Unable to load embedded static assets: %s", err)
	}
	router.Use(static.Serve("/", embedFS))
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
	fileURL := func() string {
		if isGithub {
			return fmt.Sprintf("%s/%s/%s/%s%s", GithubBaseURI, c.Param("owner"), c.Param("repo"), c.Param("commit"), c.Param("path"))
		}
		results := CleanURLSpaces(c.Param("owner"), c.Param("repo"), c.Param("commit"), c.Param("path"))
		return fmt.Sprintf("%s/%s/%s/%s/%s%s", GitLabBaseURL, results[0], results[1], "/-/raw/", results[2], results[3])

	}()
	resp, err := http.Head(fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}
	_ = resp.Body.Close()

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

	resp, err = http.Get(fileURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	defer func() {
		_ = resp.Body.Close()
	}()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err,
		})
		return
	}

	c.String(http.StatusOK, string(body[:]))
}
