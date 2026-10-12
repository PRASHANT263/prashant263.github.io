// Command build renders the site into dist/.
package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"portfolio/internal/site"
)

func main() {
	var cfg site.Config
	cfg.RegisterFlags(flag.CommandLine)
	flag.Parse()

	start := time.Now()
	if err := site.Build(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "build failed:\n%v\n", err)
		os.Exit(1)
	}
	fmt.Printf("built %s in %s\n", cfg.OutDir, time.Since(start).Round(time.Millisecond))
}
