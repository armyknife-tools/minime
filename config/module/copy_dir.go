package module

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
)

// copyDir copies the src directory contents into dst. Both directories
// should already exist.
func copyDir(dst, src string) error {
	src, err := filepath.EvalSymlinks(src)
	if err != nil {
		return err
	}

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == src {
			return nil
		}

		if strings.HasPrefix(filepath.Base(path), ".") {
			// Skip any dot files
			if info.IsDir() {
				return filepath.SkipDir
			} else {
				return nil
			}
		}

		// The "path" has the src prefixed to it. We need to join our
		// destination with the path without the src on it.
		dstPath := filepath.Join(dst, path[len(src):])

		// we don't want to try and copy the same file over itself.
		if path == dstPath {
			return nil
		}

		// We still might have the same file through a link, so check the
		// inode if we can
		if eq, err := sameInode(path, dstPath); eq {
			return nil
		} else if err != nil {
			return err
		}

		// If we have a directory, make that subdirectory, then continue
		// the walk.
		if info.IsDir() {
			if path == filepath.Join(src, dst) {
				// dst is in src; don't walk it.
				return nil
			}

			if err := os.MkdirAll(dstPath, 0755); err != nil {
				return err
			}

			return nil
		}

		// If we have a file, copy the contents.
		srcF, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcF.Close()

		dstF, err := os.Create(dstPath)
		if err != nil {
			return err
		}
		defer dstF.Close()

		if _, err := io.Copy(dstF, srcF); err != nil {
			return err
		}

		// Chmod it
		return os.Chmod(dstPath, info.Mode())
	}

	return filepath.Walk(src, walkFn)
}

// sameInode looks up the inode for paths a and b and returns if they are
// equal. On windows this will always return false.
func sameInode(a, b string) (bool, error) {
	var aIno, bIno uint64
	aStat, err := os.Stat(a)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if st, ok := aStat.Sys().(*syscall.Stat_t); ok {
		aIno = st.Ino
	}

	bStat, err := os.Stat(b)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if st, ok := bStat.Sys().(*syscall.Stat_t); ok {
		bIno = st.Ino
	}

	if aIno > 0 && aIno == bIno {
		return true, nil
	}

	return false, nil
}
