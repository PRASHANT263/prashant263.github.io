package content

import (
	"errors"
	"fmt"
	"io/fs"
	"net/url"
	"regexp"
	"strings"
)

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// Validate checks content for mistakes that would produce a broken site.
// static is the static/ directory, used to confirm referenced files exist.
// Every problem is reported, each as "file: field: message".
func Validate(s *Site, static fs.FS) error {
	var p problems
	p.profile(s.Profile, static)
	p.experience(s.Experience)
	p.skills(s.Skills)
	p.projects(s.Projects)
	p.blog(s.Blog)
	p.caseStudies(s.CaseStudies)
	p.recognition(s.Recognition)
	return errors.Join(p...)
}

type problems []error

func (p *problems) add(file, field, msg string) {
	*p = append(*p, fmt.Errorf("%s: %s: %s", file, field, msg))
}

func (p *problems) required(file, field, value string) {
	if strings.TrimSpace(value) == "" {
		p.add(file, field, "is required")
	}
}

func (p *problems) absURL(file, field, value string) {
	if value != "" && !isAbsURL(value) {
		p.add(file, field, fmt.Sprintf("%q must be an absolute http(s) URL", value))
	}
}

func (p *problems) staticFile(static fs.FS, file, field, value string) {
	name := strings.TrimPrefix(value, "/")
	if name == "" {
		return
	}
	if _, err := fs.Stat(static, name); err != nil {
		p.add(file, field, fmt.Sprintf("static/%s does not exist", name))
	}
}

func isAbsURL(s string) bool {
	u, err := url.Parse(s)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}

func (p *problems) profile(pr Profile, static fs.FS) {
	const f = "profile.yaml"
	p.required(f, "name", pr.Name)
	p.required(f, "title", pr.Title)
	p.required(f, "pitch", pr.Pitch)
	p.required(f, "photo", pr.Photo)
	p.required(f, "resume", pr.Resume)
	p.required(f, "site_url", pr.SiteURL)
	p.absURL(f, "site_url", pr.SiteURL)
	p.required(f, "email", pr.Email)
	if pr.Email != "" && !strings.Contains(pr.Email, "@") {
		p.add(f, "email", fmt.Sprintf("%q is not an email address", pr.Email))
	}
	p.staticFile(static, f, "photo", pr.Photo)
	p.staticFile(static, f, "avatar", pr.Avatar)
	p.staticFile(static, f, "resume", pr.Resume)
	p.staticFile(static, f, "resume_preview", pr.ResumePreview)
	for i, a := range pr.About {
		p.required(f, fmt.Sprintf("about[%d].title", i), a.Title)
		p.required(f, fmt.Sprintf("about[%d].text", i), a.Text)
		if !aboutIcons[a.Icon] {
			p.add(f, fmt.Sprintf("about[%d].icon", i), fmt.Sprintf("%q must be code, target or users", a.Icon))
		}
	}
	for i, s := range pr.Socials {
		field := fmt.Sprintf("socials[%d].url", i)
		p.required(f, field, s.URL)
		if !strings.HasPrefix(s.URL, "mailto:") {
			p.absURL(f, field, s.URL)
		}
	}
}

func (p *problems) experience(companies []Company) {
	const f = "experience.yaml"
	if len(companies) == 0 {
		p.add(f, "companies", "at least one company is required")
	}
	for i, c := range companies {
		p.required(f, fmt.Sprintf("[%d].name", i), c.Name)
		if len(c.Roles) == 0 {
			p.add(f, fmt.Sprintf("[%d].roles", i), "at least one role is required")
		}
		for j, r := range c.Roles {
			p.required(f, fmt.Sprintf("[%d].roles[%d].title", i, j), r.Title)
			p.required(f, fmt.Sprintf("[%d].roles[%d].start", i, j), r.Start)
		}
	}
}

func (p *problems) skills(groups []SkillGroup) {
	const f = "skills.yaml"
	for i, g := range groups {
		p.required(f, fmt.Sprintf("[%d].name", i), g.Name)
		if len(g.Items) == 0 {
			p.add(f, fmt.Sprintf("[%d].items", i), "at least one skill is required")
		}
	}
}

func (p *problems) projects(projects []Project) {
	const f = "projects.yaml"
	for i, pr := range projects {
		p.required(f, fmt.Sprintf("[%d].title", i), pr.Title)
		p.required(f, fmt.Sprintf("[%d].description", i), pr.Description)
		p.absURL(f, fmt.Sprintf("[%d].repo", i), pr.Repo)
		p.absURL(f, fmt.Sprintf("[%d].live", i), pr.Live)
	}
}

func (p *problems) blog(entries []BlogEntry) {
	const f = "blog.yaml"
	seen := map[string]bool{}
	for i, e := range entries {
		field := func(name string) string { return fmt.Sprintf("[%d].%s", i, name) }
		p.required(f, field("title"), e.Title)
		if e.Date.IsZero() {
			p.add(f, field("date"), "is required")
		}
		switch {
		case e.URL == "" && e.Slug == "":
			p.add(f, field("url"), "set either url (external link) or slug (local post)")
		case e.URL != "" && e.Slug != "":
			p.add(f, field("url"), "set only one of url or slug")
		case e.URL != "":
			p.absURL(f, field("url"), e.URL)
		default:
			if !slugPattern.MatchString(e.Slug) {
				p.add(f, field("slug"), fmt.Sprintf("%q must be lowercase letters, digits and hyphens", e.Slug))
			}
			if seen[e.Slug] {
				p.add(f, field("slug"), fmt.Sprintf("duplicate slug %q", e.Slug))
			}
			seen[e.Slug] = true
		}
		for _, tag := range e.Tags {
			if !slugPattern.MatchString(tag) {
				p.add(f, field("tags"), fmt.Sprintf("%q must be lowercase letters, digits and hyphens", tag))
			}
		}
	}
}

func (p *problems) caseStudies(studies []CaseStudy) {
	for _, cs := range studies {
		f := "case-studies/" + cs.Slug + ".md"
		if !slugPattern.MatchString(cs.Slug) {
			p.add(f, "file name", "must be lowercase letters, digits and hyphens")
		}
		p.required(f, "title", cs.Title)
		p.required(f, "summary", cs.Summary)
	}
}

var recognitionKinds = map[string]bool{"": true, "award": true, "milestone": true, "education": true}

func (p *problems) recognition(items []Recognition) {
	const f = "recognition.yaml"
	for i, r := range items {
		p.required(f, fmt.Sprintf("[%d].title", i), r.Title)
		p.required(f, fmt.Sprintf("[%d].when", i), r.When)
		if !recognitionKinds[r.Kind] {
			p.add(f, fmt.Sprintf("[%d].kind", i), fmt.Sprintf("%q must be award, milestone or education", r.Kind))
		}
	}
}

var aboutIcons = map[string]bool{"": true, "code": true, "target": true, "users": true}
