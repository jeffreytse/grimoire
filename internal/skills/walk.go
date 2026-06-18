package skills

import (
	"os"
	"path/filepath"
	"strings"
)

type Domain struct {
	Name   string
	Path   string
	Nested bool
}

type Subdomain struct {
	Domain string
	Name   string
	Path   string
}

type Skill struct {
	Domain    string
	Subdomain string
	Name      string
	Path      string
}

// IsNested returns true when the domain uses subdomain directories
// (no direct skills/ folder, or skills/ is empty).
// Defaults to true when the skills/ dir is absent so that callers iterate
// subdomains rather than assuming a flat layout for an unrecognised domain.
func IsNested(domainDir string) bool {
	skillsDir := filepath.Join(domainDir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return true
	}
	for _, e := range entries {
		if !strings.HasPrefix(e.Name(), ".") {
			return false
		}
	}
	return true
}

func ListDomains(skillsRoot string) ([]Domain, error) {
	entries, err := os.ReadDir(skillsRoot)
	if err != nil {
		return nil, err
	}
	var domains []Domain
	for _, e := range entries {
		// .claude-plugin is a Claude plugin metadata directory that lives
		// alongside skill domains but is not itself a domain.
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == ".claude-plugin" {
			continue
		}
		path := filepath.Join(skillsRoot, e.Name())
		domains = append(domains, Domain{
			Name:   e.Name(),
			Path:   path,
			Nested: IsNested(path),
		})
	}
	return domains, nil
}

func ListSubdomains(domainDir string) ([]Subdomain, error) {
	domain := filepath.Base(domainDir)
	entries, err := os.ReadDir(domainDir)
	if err != nil {
		return nil, err
	}
	var subs []Subdomain
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") || e.Name() == ".claude-plugin" {
			continue
		}
		subPath := filepath.Join(domainDir, e.Name())
		if _, err := os.Stat(filepath.Join(subPath, "skills")); err != nil {
			continue
		}
		subs = append(subs, Subdomain{
			Domain: domain,
			Name:   e.Name(),
			Path:   subPath,
		})
	}
	return subs, nil
}

func ListSkillsInDir(dir, domainName, subdomainName string) ([]Skill, error) {
	skillsDir := filepath.Join(dir, "skills")
	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		return nil, err
	}
	var skills []Skill
	for _, e := range entries {
		if strings.HasPrefix(e.Name(), ".") {
			continue
		}
		skillPath := filepath.Join(skillsDir, e.Name())
		if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
			continue
		}
		skills = append(skills, Skill{
			Domain:    domainName,
			Subdomain: subdomainName,
			Name:      e.Name(),
			Path:      skillPath,
		})
	}
	return skills, nil
}

func ListAllSkills(skillsRoot string) ([]Skill, error) {
	domains, err := ListDomains(skillsRoot)
	if err != nil {
		return nil, err
	}
	var all []Skill
	for _, d := range domains {
		if d.Nested {
			subs, err := ListSubdomains(d.Path)
			if err != nil {
				continue
			}
			for _, s := range subs {
				skills, err := ListSkillsInDir(s.Path, d.Name, s.Name)
				if err != nil {
					continue
				}
				all = append(all, skills...)
			}
		} else {
			skills, err := ListSkillsInDir(d.Path, d.Name, "")
			if err != nil {
				continue
			}
			all = append(all, skills...)
		}
	}
	return all, nil
}

// ResolveSkillPath resolves a skill reference like "engineering/development/propose-commit"
// to its absolute path in the skills root.
func ResolveSkillPath(skillsRoot, ref string) (string, error) {
	parts := strings.Split(ref, "/")
	var skillPath string
	switch len(parts) {
	case 3:
		skillPath = filepath.Join(skillsRoot, parts[0], parts[1], "skills", parts[2])
	case 2:
		skillPath = filepath.Join(skillsRoot, parts[0], "skills", parts[1])
	default:
		return "", &ErrInvalidSkillRef{Ref: ref}
	}
	if _, err := os.Stat(filepath.Join(skillPath, "SKILL.md")); err != nil {
		return "", &ErrSkillNotFound{Path: skillPath}
	}
	return skillPath, nil
}

type ErrInvalidSkillRef struct{ Ref string }
type ErrSkillNotFound struct{ Path string }

func (e *ErrInvalidSkillRef) Error() string {
	return "invalid skill reference: " + e.Ref + " (use domain/subdomain/skill or domain/skill)"
}
func (e *ErrSkillNotFound) Error() string { return "skill not found: " + e.Path }
