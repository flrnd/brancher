// Package git
package git

import "github.com/go-git/go-git/v6/plumbing"

type Repository struct {
	Path string
}

type Branch struct {
	RefName plumbing.ReferenceName
}
