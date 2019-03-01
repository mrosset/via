package via

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"os"
	gpath "path"
)

// gitref returns git branch reference
func gitref(ref string) plumbing.ReferenceName {
	return plumbing.ReferenceName(
		fmt.Sprintf("refs/heads/%s", ref),
	)
}

// Checkout git ref in dir
func Checkout(dir Path, ref string) error {
	r, err := git.PlainOpen(dir.String())
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	return w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(ref),
	})
}

// Clone remote URL into directory.
func Clone(dir Path, url string) error {
	_, err := git.PlainClone(dir.String(), false, &git.CloneOptions{
		URL:               url,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	return err
}

func references(path Path) (refs []string, err error) {
	r, err := git.PlainOpen(path.String())
	if err != nil {
		return nil, err
	}
	iter, err := r.References()
	if err != nil {
		return nil, err
	}
	fn := func(r *plumbing.Reference) error {
		refs = append(refs, string(r.Name()))
		return nil
	}
	return refs, iter.ForEach(fn)
}

// CloneBranch clone remove URL with branch to directory
func CloneBranch(dir Path, url, ref string) error {
	_, err := git.PlainClone(dir.String(), false, &git.CloneOptions{
		URL:           url,
		Progress:      os.Stdout,
		Depth:         1,
		ReferenceName: gitref(ref),
	})
	return err
}

// Branch returns the currently checked out branch for a git directory
// FIXME: this will probably fail with a detached head
func Branch(path Path) (string, error) {
	r, err := git.PlainOpen(path.String())
	if err != nil {
		return "", err
	}
	head, err := r.Head()
	if err != nil {
		return "", err
	}
	return gpath.Base(head.Name().String()), nil
}
