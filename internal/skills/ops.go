package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// InstallSkill installs a skill directory into destDir.
// If symlink is true, creates a symlink; otherwise copies the directory.
func InstallSkill(src, destDir string, symlink bool) (installed bool, err error) {
	name := filepath.Base(src)
	dest := filepath.Join(destDir, name)

	if symlink {
		if existing, err := os.Readlink(dest); err == nil {
			if existing == src {
				return false, nil // already correct
			}
			// wrong target — remove and relink
			if err := os.Remove(dest); err != nil {
				return false, fmt.Errorf("removing old symlink %s: %w", dest, err)
			}
		} else if _, err := os.Lstat(dest); err == nil {
			if _, statErr := os.Stat(dest); statErr == nil {
				// valid dir/file collision — warn and skip
				fmt.Fprintf(os.Stderr, "  warn: %s already exists at %s, skipping\n", name, dest)
				return false, nil
			}
			// broken symlink
			if err := os.Remove(dest); err != nil {
				return false, fmt.Errorf("removing broken symlink %s: %w", dest, err)
			}
		}
		if err := os.MkdirAll(destDir, 0o755); err != nil {
			return false, fmt.Errorf("creating %s: %w", destDir, err)
		}
		if err := os.Symlink(src, dest); err != nil {
			return false, fmt.Errorf("symlinking %s: %w", name, err)
		}
		return true, nil
	}

	// copy mode
	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return false, fmt.Errorf("creating %s: %w", destDir, err)
	}
	if err := copyDir(src, dest); err != nil {
		return false, fmt.Errorf("copying %s: %w", name, err)
	}
	return true, nil
}

// UninstallSkill removes a skill by name from destDir.
func UninstallSkill(name, destDir string) (removed bool, err error) {
	dest := filepath.Join(destDir, name)
	if _, err := os.Lstat(dest); err != nil {
		return false, nil // not installed
	}
	if err := os.RemoveAll(dest); err != nil {
		return false, fmt.Errorf("removing %s: %w", name, err)
	}
	return true, nil
}

// CleanBrokenSymlinks removes broken symlinks from dir. Returns count removed.
func CleanBrokenSymlinks(dir string) (int, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	count := 0
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		full := filepath.Join(dir, e.Name())
		if _, lErr := os.Lstat(full); lErr != nil {
			continue
		}
		if _, sErr := os.Stat(full); sErr != nil {
			if err := os.Remove(full); err == nil {
				count++
			}
		}
	}
	return count, nil
}

func copyDir(src, dest string) error {
	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dest, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }()
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	if _, err = io.Copy(out, in); err != nil {
		_ = out.Close()
		return err
	}
	return out.Close()
}
