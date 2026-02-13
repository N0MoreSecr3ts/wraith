package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFindingSetupUrlsGithub(t *testing.T) {
	Convey("Given a GitHub finding and session", t, func() {
		sess := &Session{
			ScanType: "github",
		}
		f := &Finding{
			RepositoryOwner: "N0MoreSecr3ts",
			RepositoryName:  "wraith",
			CommitHash:      "abcd1234",
			FilePath:        "core/findings.go",
		}

		Convey("When setupUrls is called", func() {
			f.setupUrls(sess)

			Convey("The repository, file and commit URLs should be correctly composed", func() {
				So(f.RepositoryURL, ShouldEqual, "https://github.com/N0MoreSecr3ts/wraith")
				So(f.FileURL, ShouldEqual, "https://github.com/N0MoreSecr3ts/wraith/blob/abcd1234/core/findings.go")
				So(f.CommitURL, ShouldEqual, "https://github.com/N0MoreSecr3ts/wraith/commit/abcd1234")
			})
		})
	})
}

func TestFindingSetupUrlsGitlab(t *testing.T) {
	Convey("Given a GitLab finding and session", t, func() {
		sess := &Session{
			ScanType: "gitlab",
		}
		// include spaces to exercise CleanURLSpaces behaviour
		f := &Finding{
			RepositoryOwner: "My Org",
			RepositoryName:  "My Repo",
			CommitHash:      "abcd1234",
			FilePath:        "core/findings.go",
		}

		Convey("When setupUrls is called", func() {
			f.setupUrls(sess)

			Convey("The repository, file and commit URLs should be correctly composed and cleaned", func() {
				owner, repo := CleanURLSpaces(f.RepositoryOwner, f.RepositoryName)[0], CleanURLSpaces(f.RepositoryOwner, f.RepositoryName)[1]
				So(f.RepositoryURL, ShouldEqual, "https://gitlab.com/"+owner+"/"+repo)
				So(f.FileURL, ShouldEqual, f.RepositoryURL+"/blob/abcd1234/core/findings.go")
				So(f.CommitURL, ShouldEqual, f.RepositoryURL+"/commit/abcd1234")
			})
		})
	})
}

