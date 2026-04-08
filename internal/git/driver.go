package git

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/plumbing"

	urlpkg "net/url"
)

type Driver interface {
	CreateBranch(name string) error
	CreateAndCheckoutBranch(name string) error
	DeleteBranch(name string) error
	ListLocalBranches() ([]Branch, error)
	ListRemoteBranches() ([]Branch, error)
	ListAllBranches() ([]Branch, error)
	CurrentBranch() (Branch, error)
}

type RemoteReader interface {
	OriginURL() (string, error)
}

type GoGitDriver struct {
	repo *git.Repository
}

func NewDriver() (*GoGitDriver, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		DetectDotGit: true,
	})
	if err != nil {
		return nil, err
	}

	return &GoGitDriver{
		repo: repo,
	}, nil
}

func NewRemoteReader() (RemoteReader, error) {
	return NewDriver()
}

func (d *GoGitDriver) CurrentBranch() (Branch, error) {
	ref, err := d.repo.Head()
	if err != nil {
		return Branch{}, err
	}

	return Branch{
		RefName: ref.Name(),
	}, nil
}

func (d *GoGitDriver) CreateBranch(name string) error {
	head, err := d.repo.Head()
	if err != nil {
		return err
	}

	refName := plumbing.NewBranchReferenceName(name)

	ref := plumbing.NewHashReference(refName, head.Hash())

	if err := d.repo.Storer.SetReference(ref); err != nil {
		return err
	}

	return nil
}

func (d *GoGitDriver) CreateAndCheckoutBranch(name string) error {
	worktree, err := d.repo.Worktree()
	if err != nil {
		return err
	}

	return worktree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(name),
		Create: true,
		Keep:   true,
	})
}

func (d *GoGitDriver) DeleteBranch(name string) error {
	refName := plumbing.NewBranchReferenceName(name)

	if err := d.repo.Storer.RemoveReference(refName); err != nil {
		return err
	}

	return nil
}

func (d *GoGitDriver) ListLocalBranches() ([]Branch, error) {
	return d.listBranches(func(ref *plumbing.Reference) bool {
		return ref.Name().IsBranch()
	})
}

func (d *GoGitDriver) ListRemoteBranches() ([]Branch, error) {
	return d.listBranches(func(ref *plumbing.Reference) bool {
		return ref.Name().IsRemote()
	})
}

func (d *GoGitDriver) ListAllBranches() ([]Branch, error) {
	return d.listBranches(func(ref *plumbing.Reference) bool {
		return ref.Name().IsBranch() || ref.Name().IsRemote()
	})
}

func (d *GoGitDriver) listBranches(filter func(*plumbing.Reference) bool) ([]Branch, error) {
	var branches []Branch

	refs, err := d.repo.References()
	if err != nil {
		return nil, err
	}
	defer refs.Close()
	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if filter(ref) {
			branches = append(branches, Branch{
				RefName: ref.Name(),
			})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return branches, nil
}

func (d *GoGitDriver) OriginURL() (string, error) {
	remote, err := d.repo.Remote("origin")
	if err != nil {
		return "", err
	}

	if len(remote.Config().URLs) == 0 {
		return "", fmt.Errorf("no remote URL found")
	}

	return remote.Config().URLs[0], nil
}

func ParseRemote(raw string) (owner, repo string, err error) {
	url := strings.TrimSpace(raw)

	if strings.HasPrefix(url, "git@") {
		i := strings.Index(url, ":")
		if i == -1 {
			return "", "", fmt.Errorf("invalid remote format")
		}
		url = url[i+1:]
	} else {
		u, err := urlpkg.Parse(url)
		if err != nil {
			return "", "", err
		}
		url = strings.TrimPrefix(u.Path, "/")
	}

	url = strings.TrimSuffix(url, ".git")

	parts := strings.Split(url, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid remote path")
	}

	return parts[0], parts[1], nil
}
