package content

import (
	"strings"
	"testing"
	"testing/fstest"
	"time"
)

func mustDate(s string) Date {
	t, err := time.Parse(dateLayout, s)
	if err != nil {
		panic(err)
	}
	return Date{t}
}

func validSite() *Site {
	return &Site{
		Profile: Profile{
			Name: "Ada", Title: "Engineer", Pitch: "Builds.", Email: "ada@example.com",
			Photo: "img/ada.jpg", Resume: "resume.pdf", SiteURL: "https://ada.example.com",
			Socials: []Social{
				{Label: "GitHub", Icon: "github", URL: "https://github.com/ada"},
				{Label: "Email", Icon: "email", URL: "mailto:ada@example.com"},
			},
		},
		Experience: []Company{{Name: "Acme", Roles: []Role{{Title: "Engineer", Start: "Jan 2020"}}}},
		Skills:     []SkillGroup{{Name: "Languages", Items: []string{"Go"}}},
		Projects:   []Project{{Title: "Tool", Description: "A tool.", Repo: "https://github.com/ada/tool"}},
		Blog: []BlogEntry{
			{Title: "Post", Date: mustDate("2026-01-01"), URL: "https://medium.com/@ada/post"},
			{Title: "Local", Date: mustDate("2026-01-02"), Slug: "local"},
		},
		CaseStudies: []CaseStudy{{Slug: "migration", Title: "Migration", Summary: "Moved."}},
		Recognition: []Recognition{{Title: "Star Performer", When: "2024", Kind: "award"}},
	}
}

var testStatic = fstest.MapFS{"img/ada.jpg": {}, "resume.pdf": {}}

func TestValidateValidSite(t *testing.T) {
	if err := Validate(validSite(), testStatic); err != nil {
		t.Fatal(err)
	}
}

func TestValidateProblems(t *testing.T) {
	tests := []struct {
		name   string
		mutate func(*Site)
		want   string
	}{
		{"missing name", func(s *Site) { s.Profile.Name = "" }, "profile.yaml: name: is required"},
		{"missing email", func(s *Site) { s.Profile.Email = "" }, "profile.yaml: email: is required"},
		{"relative site url", func(s *Site) { s.Profile.SiteURL = "ada.example.com" }, "profile.yaml: site_url:"},
		{"bad email", func(s *Site) { s.Profile.Email = "ada" }, "profile.yaml: email:"},
		{"bubble missing detail", func(s *Site) { s.Profile.Bubbles = []SkillBubble{{Skill: "Go"}} }, "profile.yaml: bubbles[0].detail: is required"},
		{"missing avatar file", func(s *Site) { s.Profile.Avatar = "img/none.svg" }, "profile.yaml: avatar: static/img/none.svg does not exist"},
		{"missing photo file", func(s *Site) { s.Profile.Photo = "img/none.jpg" }, "static/img/none.jpg does not exist"},
		{"bad social url", func(s *Site) { s.Profile.Socials[0].URL = "github.com/ada" }, "profile.yaml: socials[0].url:"},
		{"no companies", func(s *Site) { s.Experience = nil }, "experience.yaml: companies: at least one"},
		{"role without title", func(s *Site) { s.Experience[0].Roles[0].Title = "" }, "experience.yaml: [0].roles[0].title: is required"},
		{"empty skill group", func(s *Site) { s.Skills[0].Items = nil }, "skills.yaml: [0].items:"},
		{"bad project repo", func(s *Site) { s.Projects[0].Repo = "nope" }, "projects.yaml: [0].repo:"},
		{"blog missing date", func(s *Site) { s.Blog[0].Date = Date{} }, "blog.yaml: [0].date: is required"},
		{"blog url and slug", func(s *Site) { s.Blog[0].Slug = "x" }, "blog.yaml: [0].url: set only one of url or slug"},
		{"blog neither", func(s *Site) { s.Blog[0].URL = "" }, "blog.yaml: [0].url: set either url"},
		{"blog bad slug", func(s *Site) { s.Blog[1].Slug = "Bad Slug" }, "blog.yaml: [1].slug:"},
		{"blog duplicate slug", func(s *Site) { s.Blog[0].URL = ""; s.Blog[0].Slug = "local" }, "blog.yaml: [1].slug: duplicate"},
		{"about missing text", func(s *Site) { s.Profile.About = []AboutPoint{{Title: "What I do"}} }, "profile.yaml: about[0].text: is required"},
		{"about bad icon", func(s *Site) { s.Profile.About = []AboutPoint{{Title: "T", Text: "x", Icon: "rocket"}} }, "profile.yaml: about[0].icon:"},
		{"recognition missing when", func(s *Site) { s.Recognition[0].When = "" }, "recognition.yaml: [0].when: is required"},
		{"recognition bad kind", func(s *Site) { s.Recognition[0].Kind = "trophy" }, "recognition.yaml: [0].kind:"},
		{"case study missing summary", func(s *Site) { s.CaseStudies[0].Summary = "" }, "case-studies/migration.md: summary: is required"},
		{"blog bad tag", func(s *Site) { s.Blog[1].Tags = []string{"Good", "Bad Tag"} }, `blog.yaml: [1].tags: "Good" must be lowercase letters, digits and hyphens`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validSite()
			tt.mutate(s)
			err := Validate(s, testStatic)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want it to contain %q", err, tt.want)
			}
		})
	}
}

func TestValidateReportsEveryProblem(t *testing.T) {
	s := validSite()
	s.Profile.Name, s.Profile.Title = "", ""
	err := Validate(s, testStatic)
	for _, want := range []string{"name: is required", "title: is required"} {
		if err == nil || !strings.Contains(err.Error(), want) {
			t.Errorf("err = %v, missing %q", err, want)
		}
	}
}

func TestValidateGoogleForm(t *testing.T) {
	good := GoogleForm{Action: "https://docs.google.com/forms/d/e/1FAIpQLSf-abc_123/formResponse", Name: "entry.111", Phone: "entry.222", Email: "entry.333", Message: "entry.444"}
	tests := []struct {
		name   string
		mutate func(*GoogleForm)
		want   string
	}{
		{"valid", func(*GoogleForm) {}, ""},
		{"phone optional", func(g *GoogleForm) { g.Phone = "" }, ""},
		{"viewform url", func(g *GoogleForm) { g.Action = "https://docs.google.com/forms/d/e/abc/viewform" }, "contact_form.action:"},
		{"missing email field", func(g *GoogleForm) { g.Email = "" }, "contact_form.email: is required"},
		{"bad entry id", func(g *GoogleForm) { g.Name = "name" }, "contact_form.name:"},
		{"fields without action", func(g *GoogleForm) { g.Action = "" }, "contact_form.action: is required"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := validSite()
			s.Profile.ContactForm = good
			tt.mutate(&s.Profile.ContactForm)
			err := Validate(s, testStatic)
			if tt.want == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("err = %v, want %q", err, tt.want)
			}
		})
	}
}
