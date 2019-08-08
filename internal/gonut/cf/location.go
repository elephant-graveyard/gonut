package cf

import (
	"net/url"
	"os"

	"github.com/homeport/pina-golada/pkg/files"
	"gopkg.in/src-d/go-git.v4"
)

// HomeDir returns the HOME env key
func HomeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func cloneOrPull(location string, url string) error {
	// Gonut hasn't cloned the project before - execute clone
	if _, err := os.Stat(location); os.IsNotExist(err) {
		if _, err := git.PlainClone(location, false, &git.CloneOptions{URL: url}); err != nil {
			return err
		}

		// A local copy of the project already exists - execute pull
	} else {
		r, err := git.PlainOpen(location)
		if err != nil {
			return err
		}

		w, err := r.Worktree()
		if err != nil {
			return err
		}

		// Execute git pull in the specific worktree
		err = w.Pull(&git.PullOptions{RemoteName: "origin"})
		if err != nil && err.Error() != "already up-to-date" {
			return err
		}
	}

	return nil
}

// DownloadAppArtifact returns a in-memory directory from the provided path.
// In case of a public url, it tries to clone/pull the git repository to the
// .gonut/ directory and creates the in-memory directory from this copy.
func DownloadAppArtifact(rootURL string, relativePath string) (files.Directory, error) {
	// Split rootURL in its physical parts
	u, err := url.Parse(rootURL)
	if err != nil {
		return nil, err
	}

	// Initialize in-memory directory
	directory := files.NewRootDirectory()

	// In case of a file URI, load directory directly and return
	if u.Scheme == "file" {
		err := files.LoadFromDisk(directory, u.Path)
		if err != nil {
			return nil, err
		}

		return directory, nil
	}

	// In case of an HTTP URL, try to clone/pull and load local path into directory
	// Example path: ~/.gonut/github.com/cloudfoundry/cf-acceptance-tests/
	localPath := HomeDir() + "/.gonut/" + u.Host + u.Path
	err = cloneOrPull(localPath, rootURL)
	if err != nil {
		return nil, err
	}

	err = files.LoadFromDisk(directory, localPath+"/"+relativePath)
	if err != nil {
		return nil, err
	}

	return directory, nil
}
