package git

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	gogit "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// ErrNotGitRepo is returned (or wrapped) when a directory is not a git repository.
var ErrNotGitRepo = gogit.ErrRepositoryNotExists

// State holds the current repo state.
type State struct {
	Commit  string
	Version string
	Date    string
}

// CurrentState returns short commit hash, VERSION file content, and commit date.
func CurrentState(dir string) (State, error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return State{}, fmt.Errorf("opening repo at %s: %w", dir, err)
	}
	head, err := r.Head()
	if err != nil {
		return State{}, err
	}
	commit, err := r.CommitObject(head.Hash())
	if err != nil {
		return State{}, err
	}
	short := head.Hash().String()
	if len(short) > 7 {
		short = short[:7]
	}
	ver := readVersion(dir)
	return State{
		Commit:  short,
		Version: ver,
		Date:    commit.Author.When.Format("2006-01-02"),
	}, nil
}

// Pull fetches and merges the remote branch into the working tree.
func Pull(dir string) error {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("opening repo: %w", err)
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	err = w.Pull(&gogit.PullOptions{})
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

// PullWithForceFallback attempts a normal pull. If the remote was force-pushed
// (non-fast-forward), it warns to stderr and recovers with FetchReset.
// Pass force=true to discard local modifications instead of returning an error.
func PullWithForceFallback(dir string, force bool) error {
	err := Pull(dir)
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, gogit.ErrUnstagedChanges):
		if !force {
			return fmt.Errorf("package has local modifications; use --force to discard them and update")
		}
		fmt.Fprintf(os.Stderr, "warn: package has local modifications; discarding to update\n")
		return FetchReset(dir, force)
	case strings.Contains(err.Error(), "non-fast-forward"):
		fmt.Fprintf(os.Stderr, "warn: package was force-pushed; resetting local clone to remote HEAD\n")
		return FetchReset(dir, force)
	default:
		return err
	}
}

// FetchReset fetches from origin (with force) and hard-resets to the remote HEAD.
// Use for read-only package clones: handles force-pushed remotes without error.
// Pass force=true to discard local modifications; otherwise returns an error when dirty.
func FetchReset(dir string, force bool) error {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return fmt.Errorf("opening repo: %w", err)
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	st, err := w.Status()
	if err != nil {
		return fmt.Errorf("checking worktree status: %w", err)
	}
	if !st.IsClean() {
		if !force {
			return fmt.Errorf("package has local modifications; use --force to discard them and update")
		}
		fmt.Fprintf(os.Stderr, "warn: discarding local package modifications\n")
	}
	if err := r.Fetch(&gogit.FetchOptions{Force: true}); err != nil && !errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return fmt.Errorf("fetching: %w", err)
	}
	remoteRef, err := resolveRemoteHead(r)
	if err != nil {
		return err
	}
	return w.Reset(&gogit.ResetOptions{Commit: remoteRef.Hash(), Mode: gogit.HardReset})
}

func resolveRemoteHead(r *gogit.Repository) (*plumbing.Reference, error) {
	for _, name := range []string{"HEAD", "main", "master"} {
		if ref, err := r.Reference(plumbing.NewRemoteReferenceName("origin", name), true); err == nil {
			return ref, nil
		}
	}
	return nil, fmt.Errorf("cannot resolve remote HEAD")
}

// FetchTags fetches all tags from the remote.
func FetchTags(dir string) error {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return err
	}
	err = r.Fetch(&gogit.FetchOptions{Tags: gogit.AllTags})
	if errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return nil
	}
	return err
}

// LatestTag returns the highest semver tag in the repo.
func LatestTag(dir string) (string, error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return "", err
	}
	iter, err := r.Tags()
	if err != nil {
		return "", err
	}
	var tags []string
	if err := iter.ForEach(func(ref *plumbing.Reference) error {
		name := ref.Name().Short()
		if name != "" && (name[0] == 'v' || (name[0] >= '0' && name[0] <= '9')) {
			tags = append(tags, name)
		}
		return nil
	}); err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "", fmt.Errorf("no release tags found")
	}
	sort.Slice(tags, func(i, j int) bool {
		return compareSemver(tags[i], tags[j]) > 0
	})
	return tags[0], nil
}

// CheckoutTag checks out the given tag.
func CheckoutTag(dir, tag string) error {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return err
	}
	w, err := r.Worktree()
	if err != nil {
		return err
	}
	ref, err := r.Tag(tag)
	if err != nil {
		return fmt.Errorf("tag %s not found: %w", tag, err)
	}
	tagObj, err := r.TagObject(ref.Hash())
	var hash plumbing.Hash
	if err == nil {
		hash = tagObj.Target
	} else {
		hash = ref.Hash()
	}
	return w.Checkout(&gogit.CheckoutOptions{Hash: hash})
}

// CheckoutVersion fetches tags and checks out the given tag/branch/commit in dir.
// No-op when version is empty.
func CheckoutVersion(dir, version string) error {
	if version == "" {
		return nil
	}
	if err := FetchTags(dir); err != nil {
		return err
	}
	return CheckoutTag(dir, version)
}

// Clone clones the grimoire repo to dest.
func Clone(repoURL, dest string) error {
	_, err := gogit.PlainClone(dest, false, &gogit.CloneOptions{
		URL:   repoURL,
		Depth: 1,
	})
	return err
}

// IsUpToDate reports whether HEAD matches the remote upstream.
func IsUpToDate(dir string) (upToDate bool, local, remote State, err error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return false, State{}, State{}, err
	}
	if err := r.Fetch(&gogit.FetchOptions{Force: true}); err != nil && !errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return false, State{}, State{}, err
	}
	head, err := r.Head()
	if err != nil {
		return false, State{}, State{}, err
	}
	remoteRef, err := r.Reference(plumbing.NewRemoteReferenceName("origin", "HEAD"), true)
	if err != nil {
		// try main/master
		remoteRef, err = r.Reference(plumbing.NewRemoteReferenceName("origin", "main"), true)
		if err != nil {
			remoteRef, err = r.Reference(plumbing.NewRemoteReferenceName("origin", "master"), true)
			if err != nil {
				return false, State{}, State{}, fmt.Errorf("cannot find remote HEAD: %w", err)
			}
		}
	}
	local, _ = CurrentState(dir)
	if head.Hash() == remoteRef.Hash() {
		return true, local, local, nil
	}
	remoteCommit, err := r.CommitObject(remoteRef.Hash())
	if err != nil {
		return false, local, State{}, err
	}
	short := remoteRef.Hash().String()
	if len(short) > 7 {
		short = short[:7]
	}
	remote = State{
		Commit:  short,
		Version: remoteVersion(r, remoteCommit),
		Date:    remoteCommit.Author.When.Format("2006-01-02"),
	}
	return false, local, remote, nil
}

// PackageChanges holds the categorized artifact changes in a package since a given commit.
type PackageChanges struct {
	SkillsAdded     []string
	SkillsUpdated   []string
	ProfilesAdded   []string
	ProfilesUpdated []string
}

// HasChanges reports whether any artifact changed.
func (c *PackageChanges) HasChanges() bool {
	return len(c.SkillsAdded)+len(c.SkillsUpdated)+
		len(c.ProfilesAdded)+len(c.ProfilesUpdated) > 0
}

// HasSkillChanges reports whether any skill was added or updated.
func (c *PackageChanges) HasSkillChanges() bool {
	return len(c.SkillsAdded)+len(c.SkillsUpdated) > 0
}

// PackageChangesSince returns categorized artifact changes since oldCommit.
func PackageChangesSince(dir, oldCommit string) (PackageChanges, error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return PackageChanges{}, err
	}
	old, err := r.CommitObject(plumbing.NewHash(oldCommit))
	if err != nil {
		return PackageChanges{}, err
	}
	head, err := r.Head()
	if err != nil {
		return PackageChanges{}, err
	}
	headCommit, err := r.CommitObject(head.Hash())
	if err != nil {
		return PackageChanges{}, err
	}
	oldTree, err := old.Tree()
	if err != nil {
		return PackageChanges{}, err
	}
	headTree, err := headCommit.Tree()
	if err != nil {
		return PackageChanges{}, err
	}
	diffs, err := object.DiffTree(oldTree, headTree)
	if err != nil {
		return PackageChanges{}, err
	}

	var out PackageChanges
	for _, c := range diffs {
		from, to := c.From.Name, c.To.Name
		activePath := to
		if activePath == "" {
			activePath = from
		}
		kind, name := classifyPath(activePath)
		if kind == "" {
			continue
		}
		isAdd := from == ""
		isDel := to == ""
		_ = isDel // deletions noted but not surfaced separately yet
		switch kind {
		case "skill":
			if isAdd {
				out.SkillsAdded = append(out.SkillsAdded, name)
			} else {
				out.SkillsUpdated = append(out.SkillsUpdated, name)
			}
		case "profile":
			if isAdd {
				out.ProfilesAdded = append(out.ProfilesAdded, name)
			} else {
				out.ProfilesUpdated = append(out.ProfilesUpdated, name)
			}
		}
	}
	return out, nil
}

// classifyPath returns the artifact kind and derived display name for a package path.
func classifyPath(p string) (kind, name string) {
	switch {
	case strings.HasSuffix(p, "SKILL.md"):
		return "skill", skillPathToName(p)
	case strings.HasPrefix(p, "profiles/") && strings.HasSuffix(p, ".toml"):
		return "profile", strings.TrimSuffix(strings.TrimPrefix(p, "profiles/"), ".toml")
	}
	return "", ""
}

// CommitsSince returns one-line summaries ("abc1234 subject") for commits reachable
// from HEAD but not from oldCommit, newest first. Used when no artifact changed.
func CommitsSince(dir, oldCommit string) ([]string, error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return nil, err
	}
	head, err := r.Head()
	if err != nil {
		return nil, err
	}
	oldHash := plumbing.NewHash(oldCommit)
	iter, err := r.Log(&gogit.LogOptions{From: head.Hash()})
	if err != nil {
		return nil, err
	}
	var lines []string
	for {
		c, err := iter.Next()
		if err != nil {
			break
		}
		if c.Hash == oldHash {
			break
		}
		short := c.Hash.String()[:7]
		subject := strings.SplitN(strings.TrimSpace(c.Message), "\n", 2)[0]
		lines = append(lines, short+" "+subject)
	}
	return lines, nil
}

// skillPathToName converts "skills/engineering/development/apply-solid/SKILL.md"
// to "engineering/development/apply-solid".
func skillPathToName(p string) string {
	p = strings.TrimPrefix(p, "skills/")
	p = strings.TrimSuffix(p, "/SKILL.md")
	return p
}

func readVersion(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "VERSION"))
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(string(data))
}

func remoteVersion(_ *gogit.Repository, commit *object.Commit) string {
	tree, err := commit.Tree()
	if err != nil {
		return "unknown"
	}
	f, err := tree.File("VERSION")
	if err != nil {
		return "unknown"
	}
	content, err := f.Contents()
	if err != nil {
		return "unknown"
	}
	return strings.TrimSpace(content)
}

// compareSemver returns 1 if a > b, -1 if a < b, 0 if equal (simple string semver).
func compareSemver(a, b string) int {
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")
	aParts := strings.Split(a, ".")
	bParts := strings.Split(b, ".")
	maxParts := len(aParts)
	if len(bParts) > maxParts {
		maxParts = len(bParts)
	}
	for i := 0; i < maxParts; i++ {
		av, bv := 0, 0
		if i < len(aParts) {
			_, _ = fmt.Sscanf(aParts[i], "%d", &av)
		}
		if i < len(bParts) {
			_, _ = fmt.Sscanf(bParts[i], "%d", &bv)
		}
		if av > bv {
			return 1
		}
		if av < bv {
			return -1
		}
	}
	return 0
}

// TagState returns the State for a given tag.
func TagState(dir, tag string) (State, error) {
	r, err := gogit.PlainOpen(dir)
	if err != nil {
		return State{}, err
	}
	ref, err := r.Tag(tag)
	if err != nil {
		return State{}, err
	}
	tagObj, err := r.TagObject(ref.Hash())
	var hash plumbing.Hash
	if err == nil {
		hash = tagObj.Target
	} else {
		hash = ref.Hash()
	}
	commit, err := r.CommitObject(hash)
	if err != nil {
		return State{}, err
	}
	short := hash.String()
	if len(short) > 7 {
		short = short[:7]
	}
	_ = time.Now() // ensure time is imported
	return State{
		Commit:  short,
		Version: strings.TrimPrefix(tag, "v"),
		Date:    commit.Author.When.Format("2006-01-02"),
	}, nil
}
