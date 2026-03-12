package git

type Driver interface {
	CreateBranch(name string) error
	DeleteBranch(name string) error
	ListBranches() ([]Branch, error)
	CurrentBranch() (Branch, error)
}

type GoGitDriver struct {
	repo *Repository
}
