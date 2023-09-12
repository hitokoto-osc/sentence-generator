// Package utils is intended to provide some useful functions
package utils

import (
	"fmt"
	"github.com/cockroachdb/errors"
	"os"
	"time"

	"github.com/go-git/go-git/v5"
	gitConfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/hitokoto-osc/sentence-generator/config"
	"github.com/hitokoto-osc/sentence-generator/logging"
)

// GetGitAuth will return transport.AuthMethod interface by driver that configure in config file
func GetGitAuth() (transport.AuthMethod, error) {
	defer logging.Logger.Sync()
	if config.Git.Driver == "ssh" {
		auth, err := ssh.NewPublicKeys(config.Git.SSH.User, []byte(config.Git.SSH.PrivateKey), config.Git.SSH.Password)
		if err != nil {
			logging.Logger.Error(config.Git.SSH.PrivateKey)
			return nil, errors.WithStack(errors.WithMessage(err, "SSHKeyParseError"))
		}
		return auth, nil
	} else if config.Git.Driver == "http" {
		return &http.BasicAuth{
			Username: config.Git.HTTP.User,
			Password: config.Git.HTTP.Password,
		}, nil
	}
	return nil, errors.WithStack(errors.New("unsupported auth mode: " + config.Git.Driver))
}

// SyncRepository will force sync local repository correspond to remote repository
func SyncRepository() error {
	defer logging.Logger.Sync()
	logging.Logger.Info("Open local repository...")
	repository, err := git.PlainOpen(config.Core.Workdir)
	if err != nil {
		return err
	}
	logging.Logger.Info("Fetch from remote repository...")
	auth, err := GetGitAuth()
	if err != nil {
		return err
	}
	if err = repository.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}); err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			logging.Logger.Info("Local repository is up-to-date.")
			return nil
		}
		return err
	}
	logging.Logger.Info("Reset local records to latest remote version...")
	latestCommit, err := repository.Reference(plumbing.ReferenceName("refs/remotes/origin/"+config.Git.Branch), true)
	if err != nil {
		return err
	}
	WorkTree, err := repository.Worktree()
	if err != nil {
		return err
	}
	return WorkTree.Reset(&git.ResetOptions{
		Commit: latestCommit.Hash(),
		Mode:   git.HardReset,
	})
}

// CommitRepository will commit all local changes
func CommitRepository() error {
	defer logging.Logger.Sync()
	logging.Logger.Info("Open local repository...")
	repository, err := git.PlainOpen(config.Core.Workdir)
	if err != nil {
		return err
	}
	logging.Logger.Info("Do commit work...")
	version, err := GetBundleVersion()
	if err != nil {
		return err
	}
	workTree, err := repository.Worktree()
	if err != nil {
		return err
	}
	_, err = workTree.Add(".")
	if err != nil {
		return err
	}
	status, err := workTree.Status()
	if err != nil {
		return err
	}
	logging.Logger.Info(status.String())
	_, err = workTree.Commit(
		fmt.Sprintf("build: v%s", version),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  config.Git.Name,
				Email: config.Git.Email,
				When:  time.Now(),
			},
		},
	)
	return err
}

// ReleaseRepository will create a git tag
func ReleaseRepository() error {
	defer logging.Logger.Sync()
	logging.Logger.Info("Open local repository...")
	repository, err := git.PlainOpen(config.Core.Workdir)
	if err != nil {
		return err
	}
	logging.Logger.Info("Do release work...")
	version, err := GetBundleVersion()
	if err != nil {
		return err
	}
	ref, err := repository.Head()
	if err != nil {
		return err
	}
	_, err = repository.CreateTag(fmt.Sprintf("v%s", version), ref.Hash(), &git.CreateTagOptions{
		Tagger: &object.Signature{
			Name:  config.Git.Name,
			Email: config.Git.Email,
			When:  time.Now(),
		},
		Message: fmt.Sprintf("release v%s", version),
	})
	return err
}

// Push will push local changes(includes tags) to remote repository
func Push() error {
	defer logging.Logger.Sync()
	logging.Logger.Info("Open local repository...")
	repository, err := git.PlainOpen(config.Core.Workdir)
	if err != nil {
		return err
	}
	logging.Logger.Debug("Get Auth instance...")
	auth, err := GetGitAuth()
	if err != nil {
		return err
	}

	logging.Logger.Info("Pushing commits...")
	if err := repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       auth,
		Progress:   os.Stdout,
	}); err != nil {
		return err
	}

	logging.Logger.Info("Pushing Tags...")
	return repository.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		RefSpecs:   []gitConfig.RefSpec{"refs/tags/*:refs/tags/*"},
		Auth:       auth,
	})
}
