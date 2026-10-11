package content

// Site is everything the generator knows about, loaded from content/.
type Site struct {
	Profile     Profile
	Experience  []Company
	Skills      []SkillGroup
	Projects    []Project
	Blog        []BlogEntry // newest first
	CaseStudies []CaseStudy // by Order, then Slug
}

type Profile struct {
	Name          string   `yaml:"name"`
	Title         string   `yaml:"title"`
	Pitch         string   `yaml:"pitch"`
	About         []string `yaml:"about"`
	Location      string   `yaml:"location"`
	Email         string   `yaml:"email"`
	Photo         string   `yaml:"photo"`          // path under static/
	Resume        string   `yaml:"resume"`         // path under static/
	ResumePreview string   `yaml:"resume_preview"` // optional, path under static/
	ResumeUpdated Date     `yaml:"resume_updated"`
	Web3FormsKey  string   `yaml:"web3forms_key"`
	SiteURL       string   `yaml:"site_url"`
	Facts         []Fact   `yaml:"facts"`
	Socials       []Social `yaml:"socials"`
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
