package module

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
)

// GitGetter is a Getter implementation that will download a module from
// a git repository.
type GitGetter struct{}

func (g *GitGetter) Get(dst string, u *url.URL) error {
	if _, err := exec.LookPath("git"); err != nil {
		return fmt.Errorf("git must be available and on the PATH")
	}

	// Extract some query parameters we use
	q := u.Query()
	tag := q.Get("tag")
	q.Del("tag")
	u.RawQuery = q.Encode()

	// First: clone or update the repository
	_, err := os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		err = g.update(dst, u)
	} else {
		err = g.clone(dst, u)
	}
	if err != nil {
		return err
	}

	// Next: check out the proper tag/branch if it is specified, and checkout
	if tag == "" {
		return nil
	}

	return g.checkout(dst, tag)
}

func (g *GitGetter) checkout(dst string, ref string) error {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Dir = dst
	return getRunCommand(cmd)
}

func (g *GitGetter) clone(dst string, u *url.URL) error {
	cmd := exec.Command("git", "clone", u.String(), dst)
	return getRunCommand(cmd)
}

func (g *GitGetter) update(dst string, u *url.URL) error {
	cmd := exec.Command("git", "pull", "--ff-only")
	cmd.Dir = dst
	return getRunCommand(cmd)
}
