package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/N0MoreSecr3ts/wraith/core"

	ot "github.com/otiai10/copy"
	"github.com/spf13/cobra"
	whilp "github.com/whilp/git-urls"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

var signatureVersion string

// cleanInput will ensure that any user supplied git url is in the proper format
func cleanInput(u string) string {
	_, err := whilp.Parse(u)

	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	return u

}

// executeTests will run any tests associated with the expressions
// TODO deal with this
func executeTests(dir string) bool {

	// run some tests here and return a true/false depending on the outcome
	return true
}

// fetchSignatures will download the signatures from a remote location to a temp location
func fetchSignatures(sess *core.Session) string {

	// TODO if this is not set then pull from the stock place, that should be the default url set in the session
	rURL := wraithConfig.GetString("signatures-url")

	// set the remote url that we will fetch
	// TODO need to look into this more
	remoteURL := cleanInput(rURL)

	// TODO document this
	dir, err := ioutil.TempDir("", "wraith")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	// for now we only pull from a given version at some point we can look at pulling the latest
	// TODO be able to pass in a commit or version string
	repo, err := git.PlainClone(dir, false, &git.CloneOptions{
		URL:           remoteURL,
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", "master")),
		SingleBranch:  true,
		Tags:          git.AllTags,
	})
	if err != nil {
		err1 := os.RemoveAll(dir) // TODO fix this error thing
		if err1 != nil {
			sess.Out.Error(err1.Error())
		}

		sess.Out.Error(err.Error())
	}

	// TODO give a valid error if the version is not REMOVE ME
	if signatureVersion != "" {

		// Get the working tree so we can change refs
		// TODO figure this out REMOVE ME
		tree, err := repo.Worktree()
		if err != nil {
			sess.Out.Error(err.Error())
		}

		// Set the tag to the signatures version that we want to use
		// TODO fix this REMOVE ME
		tagName := string(signatureVersion)

		// Checkout our tag
		// TODO way are we using a tag here is we only checkout master
		// TODO fix this
		err = tree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName("refs/tags/" + tagName),
		})
		if err != nil {
			fmt.Println("Requested version not available. Please enter a valid version")
			os.Exit(2)
		}
	}
	return dir
}

// updateSignatures will install the new signatures into the specified location, changing the name of the previous set
func updateSignatures(rRepo string, sess *core.Session) bool {

	// create a temp directory to hold the signatures we pull
	// TODO put this in /tmp via a real library
	tempSignaturesDir := rRepo + "/signatures"

	// final resting place for the signatures
	rPath := wraithConfig.GetString("signatures-path")

	// ensure we have the proper home directory
	rPath = core.SetHomeDir(rPath, sess)

	// if the signatures path does not exist then we create it
	if !core.PathExists(rPath, sess) {

		err := os.MkdirAll(rPath, 0700)
		if err != nil {
			sess.Out.Error(err.Error())
		}
	}

	// if we want to test the signatures before we install them
	// TODO need to implement something here
	if wraithConfig.GetBool("test-signatures") {

		// if the tests pass then we install the signatures
		if executeTests(rRepo) {

			// copy the files from the temp directory to the signatures directory
			if err := ot.Copy(tempSignaturesDir, rPath); err != nil {
				sess.Out.Error(err.Error())
				return false
			}

			// get all the files in the signatures directory
			files, err := ioutil.ReadDir(rPath)
			if err != nil {
				sess.Out.Error(err.Error())
				return false
			}

			// set them to the current user and the proper permissions
			for _, f := range files {
				if err := os.Chmod(rPath+"/"+f.Name(), 0644); err != nil {
					sess.Out.Error(err.Error())
					return false
				}
			}
			err = os.RemoveAll(rRepo)
			if err != nil {
				sess.Out.Error(err.Error())
			}
			return true

		}
		err := os.RemoveAll(rRepo)
		if err != nil {
			sess.Out.Error(err.Error())
		}
		return false

	}

	// copy the files from the temp directory to the signatures directory
	if err := ot.Copy(tempSignaturesDir, rPath); err != nil {
		sess.Out.Error(err.Error())
		return false
	}

	// get all the files in the signatures directory
	files, err := ioutil.ReadDir(rPath)
	if err != nil {
		sess.Out.Error(err.Error())
		return false
	}

	// set them to the current user and the proper permissions
	// TODO ensure these are .yaml somehow
	for _, f := range files {
		sFileExt := filepath.Ext(rPath + "/" + f.Name())
		if sFileExt == "yml" || sFileExt == "yaml" {
			if err := os.Chmod(rPath+"/"+f.Name(), 0644); err != nil {
				sess.Out.Error(err.Error())
				return false
			}
		}
	}
	// TODO why is the commented out
	// TODO Cleanup after ourselves and remove any temp garbage
	// os.RemoveAll(tempSignaturesDir)
	return true
}

// updateSignaturesCmd represents the updateSignatures command
var updateSignaturesCmd = &cobra.Command{
	Use:   "updateSignatures",
	Short: "Update the signatures to the latest version available",
	Long:  "Update the signatures to the latest version available",
	Run: func(cmd *cobra.Command, args []string) {

		scanType := "updateSignatures"

		sess := core.NewSession(wraithConfig, scanType)

		// get the signatures version or if blank, set it to latest
		// TODO this should be in the default values from the session
		if wraithConfig.GetString("signatures-path") != "" {

			signatureVersion = wraithConfig.GetString("signatures-version")
		} else {
			signatureVersion = "latest"
		}

		// fetch the signatures from the remote location
		rRepo := fetchSignatures(sess)

		// install the signatures
		if updateSignatures(rRepo, sess) {
			// TODO set this in the session so we have a single location for everything
			fmt.Printf("The signatures have been successfully updated at: %s\n", wraithConfig.GetString("signatures-path"))
		} else {
			sess.Out.Warn("The signatures were not updated")
		}
	},
}

func init() {
	rootCmd.AddCommand(updateSignaturesCmd)

	updateSignaturesCmd.Flags().Bool("test-signatures", false, "run any tests associated with the signatures and display the output")
	updateSignaturesCmd.Flags().String("signatures-path", "$HOME/.wraith/signatures/", "path where the signatures will be installed")
	updateSignaturesCmd.Flags().String("signatures-url", "https://github.com/N0MoreSecr3ts/wraith-signatures", "url where the signatures can be found")
	updateSignaturesCmd.Flags().String("signatures-version", "", "specific version of the signatures to install")

	err := wraithConfig.BindPFlag("test-signatures", updateSignaturesCmd.Flags().Lookup("test-signatures"))
	err = wraithConfig.BindPFlag("signatures-path", updateSignaturesCmd.Flags().Lookup("signatures-path"))
	err = wraithConfig.BindPFlag("signatures-url", updateSignaturesCmd.Flags().Lookup("signatures-url"))
	err = wraithConfig.BindPFlag("signatures-version", updateSignaturesCmd.Flags().Lookup("signatures-version"))

	if err != nil {
		fmt.Printf("There was an error binding a flag: %s\n", err.Error())
	}
}
