package content

// Site is everything the generator knows about, loaded from content/.
type Site struct {
	Profile     Profile
	Experience  []Company
	Skills      []SkillGroup
	Projects    []Project
	Blog        []BlogEntry // newest first
	CaseStudies []CaseStudy // by Order, then Slug
	Recognition []Recognition
}

type Profile struct {
	Name          string       `yaml:"name"`
	Title         string       `yaml:"title"`
	Pitch         string       `yaml:"pitch"`
	Greetings     []string     `yaml:"greetings"` // optional; the hero types each one in turn
	Status        string       `yaml:"status"`    // optional pill above the name, e.g. "Currently building at KreditBee"
	Expertise     []string     `yaml:"expertise"` // optional footer list
	About         []AboutPoint `yaml:"about"`
	Location      string       `yaml:"location"`
	Email         string       `yaml:"email"`
	Photo         string       `yaml:"photo"`          // path under static/
	Avatar        string       `yaml:"avatar"`         // optional round monogram left of the name, path under static/
	Resume        string       `yaml:"resume"`         // path under static/
	ResumePreview string       `yaml:"resume_preview"` // optional, path under static/
	ResumeUpdated Date         `yaml:"resume_updated"`
	ContactForm   GoogleForm   `yaml:"contact_form"`
	SiteURL       string       `yaml:"site_url"`
	Facts         []Fact       `yaml:"facts"`
	Socials       []Social     `yaml:"socials"`
}

type Fact struct {
	Value string `yaml:"value"`
	Label string `yaml:"label"`
}

type Social struct {
	Label string `yaml:"label"`
	Icon  string `yaml:"icon"` // github | linkedin | email
	URL   string `yaml:"url"`
}

type Company struct {
	Name     string `yaml:"name"`
	Location string `yaml:"location"`
	Roles    []Role `yaml:"roles"`
}

type Role struct {
	Title      string   `yaml:"title"`
	Start      string   `yaml:"start"`
	End        string   `yaml:"end"` // empty means Present
	Highlights []string `yaml:"highlights"`
	Bullets    []string `yaml:"bullets"`
}

type SkillGroup struct {
	Name  string   `yaml:"name"`
	Items []string `yaml:"items"`
}

type Project struct {
	Title       string   `yaml:"title"`
	Description string   `yaml:"description"`
	Tech        []string `yaml:"tech"`
	Repo        string   `yaml:"repo"`
	Live        string   `yaml:"live"`
}

// BlogEntry is either an external link (URL set) or a local post (Slug set,
// body read from content/posts/<slug>.md).
type BlogEntry struct {
	Title   string   `yaml:"title"`
	Date    Date     `yaml:"date"`
	URL     string   `yaml:"url"`
	Slug    string   `yaml:"slug"`
	Source  string   `yaml:"source"`
	Tags    []string `yaml:"tags"`
	Summary string   `yaml:"summary"`
	Body    string   `yaml:"-"` // Markdown, local posts only
}

func (e BlogEntry) IsExternal() bool { return e.URL != "" }

// Href is where the entry's link points.
func (e BlogEntry) Href() string {
	if e.IsExternal() {
		return e.URL
	}
	return "/blog/" + e.Slug + "/"
}

type CaseStudy struct {
	Slug    string   `yaml:"-"` // from the file name
	Title   string   `yaml:"title"`
	Summary string   `yaml:"summary"`
	Order   int      `yaml:"order"`
	Tags    []string `yaml:"tags"`
	Metrics []string `yaml:"metrics"`
	Body    string   `yaml:"-"` // Markdown
}

// Recognition is an award, career milestone or degree shown in the
// Recognition section.
type Recognition struct {
	Title  string `yaml:"title"`
	When   string `yaml:"when"` // free text, e.g. "Nov 2023 · Aug 2024 · May 2025"
	Detail string `yaml:"detail"`
	Kind   string `yaml:"kind"` // award | milestone | education; picks the icon
}

// AboutPoint is one titled card in the About section.
type AboutPoint struct {
	Title string `yaml:"title"`
	Icon  string `yaml:"icon"` // code | target | users; optional
	Text  string `yaml:"text"`
}

// GoogleForm sends the contact form straight to a Google Form, whose
// responses land in its linked Google Sheet. Each field holds the
// question's "entry.<number>" ID from the form's pre-filled link.
type GoogleForm struct {
	Action  string `yaml:"action"` // https://docs.google.com/forms/d/e/<form-id>/formResponse
	Name    string `yaml:"name"`
	Phone   string `yaml:"phone"` // optional question
	Email   string `yaml:"email"`
	Message string `yaml:"message"`
}

// Enabled reports whether a form is configured; otherwise the page shows
// the email fallback.
func (g GoogleForm) Enabled() bool { return g.Action != "" }
