package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// managedMarker is written inside every copy-mode skill install so uninstall
// can safely distinguish grimoire-managed copies from user-created directories.
const managedMarker = ".grimoire-managed"

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

	// copy mode — guard against overwriting an unmanaged directory
	if info, lErr := os.Lstat(dest); lErr == nil {
		isSymlink := info.Mode()&os.ModeSymlink != 0
		if !isSymlink {
			if _, mErr := os.Stat(filepath.Join(dest, managedMarker)); mErr != nil {
				fmt.Fprintf(os.Stderr, "  warn: %s already exists and is not managed by grimoire, skipping\n", name)
				return false, nil
			}
		}
	}

	if err := os.MkdirAll(destDir, 0o755); err != nil {
		return false, fmt.Errorf("creating %s: %w", destDir, err)
	}
	if err := copyDir(src, dest); err != nil {
		return false, fmt.Errorf("copying %s: %w", name, err)
	}
	// write marker so uninstall can confirm this copy is grimoire-managed
	if err := os.WriteFile(filepath.Join(dest, managedMarker), []byte(src+"\n"), 0o644); err != nil {
		return false, fmt.Errorf("writing marker in %s: %w", name, err)
	}
	return true, nil
}

// UninstallSkill removes a skill by name from destDir.
// Symlinks are always safe to remove. Real directories require the grimoire
// managed marker — without it the directory is not touched.
func UninstallSkill(name, destDir string) (removed bool, err error) {
	dest := filepath.Join(destDir, name)
	info, err := os.Lstat(dest)
	if err != nil {
		return false, nil // not installed
	}
	if info.Mode()&os.ModeSymlink != 0 {
		// symlink: remove the link only, never the target
		if err := os.Remove(dest); err != nil {
			return false, fmt.Errorf("removing %s: %w", name, err)
		}
		return true, nil
	}
	// real directory: require grimoire managed marker before deleting
	if _, mErr := os.Stat(filepath.Join(dest, managedMarker)); mErr != nil {
		fmt.Fprintf(os.Stderr, "  warn: %s has no grimoire marker — skipping (not managed by grimoire)\n", name)
		return false, nil
	}
	if err := os.RemoveAll(dest); err != nil {
		return false, fmt.Errorf("removing %s: %w", name, err)
	}
	return true, nil
}

// CleanBrokenSymlinks removes stale grimoire-managed entries from dir:
//   - broken symlinks (Lstat ok, Stat fails)
//   - copy-mode directories whose source path (recorded in the managed marker)
//     no longer exists on disk
//
// Returns count removed.
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
		info, lErr := os.Lstat(full)
		if lErr != nil {
			continue
		}
		if info.Mode()&os.ModeSymlink != 0 {
			// broken symlink: Lstat ok but Stat fails
			if _, sErr := os.Stat(full); sErr != nil {
				if os.Remove(full) == nil {
					count++
				}
			}
			continue
		}
		// real directory: stale if it has a grimoire marker whose source is gone
		markerPath := filepath.Join(full, managedMarker)
		src, mErr := os.ReadFile(markerPath)
		if mErr != nil {
			continue // no marker — not managed by grimoire, leave it alone
		}
		sourcePath := strings.TrimSpace(string(src))
		if _, sErr := os.Stat(sourcePath); sErr != nil {
			// source gone — remove the stale copy
			if os.RemoveAll(full) == nil {
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
