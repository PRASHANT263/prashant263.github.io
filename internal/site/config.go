package site

import "flag"

// Config says where Build reads from and writes to.
type Config struct {
	ContentDir   string
	TemplatesDir string
	StaticDir    string
	OutDir       string
}

// RegisterFlags adds -content, -templates, -static and -out flags defaulting
// to the repo layout.
func (c *Config) RegisterFlags(flags *flag.FlagSet) {
	flags.StringVar(&c.ContentDir, "content", "content", "content directory")
	flags.StringVar(&c.TemplatesDir, "templates", "templates", "templates directory")
	flags.StringVar(&c.StaticDir, "static", "static", "static files directory")
	flags.StringVar(&c.OutDir, "out", "dist", "output directory (replaced on every build)")
}
