# Part 4 — Real content, deploy docs, verification

### Task 11: Real content + OG image

**Files:** Create `content/{profile,experience,skills,projects,blog}.yaml`, `content/case-studies/{ddd-migration,nbfc-co-lending,loan-documents}.md`, `assets/og-card.html`, `static/img/og.png`

Source of truth is `Prashants_a25.pdf`. Case studies use only facts on the resume: no internal service names, partner API details or unpublished metrics.

- [ ] **Step 1: `content/profile.yaml`**

```yaml
name: Prashant
title: Senior Software Engineer · Backend Systems
pitch: I build reliable, low-latency backend systems for high-scale fintech lending, in Go, on AWS, with Domain-Driven Design.
location: Bengaluru, India
email: prashantab8505@gmail.com
site_url: https://prashant263.github.io
photo: img/profile.svg          # change to img/profile.jpg after adding your photo
resume: resume.pdf
resume_preview: img/resume-preview.png
resume_updated: 2026-08-31
web3forms_key: ""               # free key from https://web3forms.com (enables the contact form)
about:
  - "I'm a Senior Software Engineer at KreditBee with 4+ years of experience building and scaling the backend systems behind high-growth digital lending."
  - "My work sits where business, regulation and engineering meet: integrating NBFC co-lending partners, generating RBI-compliant loan documents, and moving a tightly coupled monolith to a Domain-Driven Design architecture with clear service boundaries."
  - "I care about production reliability, clean architecture and helping engineers grow. I mentor five engineers through code reviews, production debugging and design discussions, and led my team's adoption of AI-assisted development with Cursor and Claude Code."
facts:
  - { value: "4+ yrs", label: "Backend engineering" }
  - { value: "3×", label: "Star Performer award" }
  - { value: "+20%", label: "Co-lending loan volume" }
  - { value: "5", label: "Engineers mentored" }
socials:
  - { label: GitHub, icon: github, url: "https://github.com/prashant263" }
  - { label: LinkedIn, icon: linkedin, url: "https://www.linkedin.com/" }   # replace with your profile URL
  - { label: Email, icon: email, url: "mailto:prashantab8505@gmail.com" }
```

- [ ] **Step 2: `content/experience.yaml`**

```yaml
- name: KreditBee
  location: Bengaluru
  roles:
    - title: Senior Software Engineer
      start: Mar 2026
      highlights: ["DDD migration", "Mentoring 5 engineers", "AI-assisted development"]
      bullets:
        - Led the migration of lending systems from a tightly coupled monolith to a Domain-Driven Design model, introducing dedicated domain, service and gateway layers to establish clear service boundaries and improve maintainability and feature delivery.
        - Mentoring 5 engineers through code reviews, production debugging and fund-flow design discussions, reinforcing engineering best practices and improving code quality.
        - Spearheaded team-wide adoption of AI-assisted development tools, including Cursor and Claude Code, to accelerate implementation and reduce developer toil.
        - Drove architectural improvements across critical lending workflows supporting multiple lending products and partner integrations, improving maintainability and production reliability.
    - title: Software Engineer
      start: Dec 2022
      end: Mar 2026
      highlights: ["+20% co-lending volume", "4 NBFC partners", "10-language documents", "−30% developer toil"]
      bullets:
        - Increased co-lending loan volume by 20% through end-to-end integration of 4 NBFC partners (Tata Capital, OFL, ABCL and IFMR), expanding the platform's lending distribution network.
        - Designed and built a scalable, low-latency Loan Document Generation service using HTML templates and a headless Chrome PDF pipeline, supporting RBI-compliant loan agreements across 10 languages.
        - Designed and shipped the Multi-Lender Display feature, owning the server-side business logic that lets users compare and select their preferred lender within a single loan journey.
        - Reduced developer toil by 30% through backend automation, streamlined API workflows and targeted system optimisations.
        - Handled production debugging and critical fund-flow issues, identifying root causes across backend services and improving reliability of lending workflows.
        - Collaborated with product, finance, compliance and operations teams to ship RBI-compliant fintech features and circulars on schedule.
    - title: Software Development Engineer Intern
      start: Apr 2022
      end: Nov 2022
      highlights: ["Go · MySQL · AWS", "Fund-flow observability"]
      bullets:
        - Ramped up on the production backend stack (Go, MySQL and AWS) alongside ReactJS, then independently owned fund-related bug fixes and product tasks.
        - Shipped lender website link integrations for all co-lending partners, and added fund assignment events that surface fund rejections and the specific rule that failed, improving debugging of fund-flow issues.
```

- [ ] **Step 3: `content/skills.yaml`**

```yaml
- { name: Languages, items: [Go, Python, SQL] }
- { name: Concurrency & Systems, items: [Goroutines, Channels, Context, Synchronization, Gin, REST APIs] }
- { name: Backend & Architecture, items: [Microservices, Domain-Driven Design, Distributed Systems, Event-Driven Architecture, API Design, Design Patterns] }
- { name: Cloud & DevOps, items: [AWS EC2, Lambda, API Gateway, S3, SNS/SQS, CloudWatch, CI/CD] }
- { name: Data & Caching, items: [MySQL, Redis] }
- { name: Observability & Tools, items: [Grafana, Monitoring, Git/GitHub, Postman, Cursor, Claude Code] }
```

- [ ] **Step 4: `content/projects.yaml`**

```yaml
- title: This portfolio
  description: A custom static site generator in Go. Content lives in YAML and Markdown, is validated at build time, rendered with html/template and goldmark, link-checked, and deployed to GitHub Pages by GitHub Actions.
  tech: [Go, html/template, goldmark, CSS, GitHub Actions]
```

- [ ] **Step 5: `content/blog.yaml`**

```yaml
# Add a blog link: copy the example below, remove the leading "# ", fill it in, save.
# Newest entries appear first automatically.
#
# - title: "How we migrated a lending monolith to DDD"
#   date: 2026-10-01                      # YYYY-MM-DD
#   url: https://medium.com/@you/your-post
#   source: Medium                        # badge text: Medium, Dev.to, Hashnode, LinkedIn...
#   tags: [go, ddd]
#   summary: "One or two sentences about the post."
```

- [ ] **Step 6: `content/case-studies/ddd-migration.md`**

```markdown
---
title: From Monolith to Domain-Driven Design
summary: Restructuring tightly coupled lending systems into domain, service and gateway layers with clear boundaries.
order: 1
tags: [go, ddd, architecture]
metrics: ["Clear service boundaries", "Faster feature delivery", "Better maintainability"]
---
## Problem

The lending platform had grown as a tightly coupled monolith. Business rules for multiple lending products and partner integrations lived alongside database access and third-party calls, so a change in one flow could ripple into others. New features meant touching code far outside the area being changed.

## Approach

I led the migration to a Domain-Driven Design model built around three layers:

- **Domain layer**: the core lending concepts and rules, in plain Go, with no knowledge of databases, HTTP or partners.
- **Service layer**: use cases that orchestrate domain logic for a specific product flow.
- **Gateway layer**: adapters for everything external, such as databases, partner systems, queues and notifications.

The migration was incremental rather than a big-bang rewrite. Flows moved one at a time behind stable interfaces, so production traffic kept flowing throughout.

## Architecture

<svg class="diagram" viewBox="0 0 770 300" role="img" aria-label="Handlers call the service layer, which uses the domain layer and the gateway layer; gateways talk to the database, partner systems and queues">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="285" y="16" width="200" height="48" rx="10"/><text x="385" y="40">API handlers</text>
<rect class="box" x="285" y="110" width="200" height="56" rx="10"/><text x="385" y="130">Service layer</text><text class="label" x="385" y="150">product use cases</text>
<rect class="box" x="20" y="110" width="200" height="56" rx="10"/><text x="120" y="130">Domain layer</text><text class="label" x="120" y="150">lending rules, pure Go</text>
<rect class="box alt" x="550" y="110" width="200" height="56" rx="10"/><text x="650" y="130">Gateway layer</text><text class="label" x="650" y="150">adapters</text>
<rect class="box alt" x="420" y="232" width="105" height="44" rx="10"/><text x="472" y="254">MySQL</text>
<rect class="box alt" x="535" y="232" width="110" height="44" rx="10"/><text x="590" y="254">Partners</text>
<rect class="box alt" x="655" y="232" width="105" height="44" rx="10"/><text x="707" y="254">Queues</text>
<line class="edge" x1="385" y1="64" x2="385" y2="106"/>
<line class="edge" x1="285" y1="138" x2="224" y2="138"/>
<line class="edge" x1="485" y1="138" x2="546" y2="138"/>
<line class="edge" x1="630" y1="166" x2="480" y2="228"/>
<line class="edge" x1="645" y1="166" x2="595" y2="228"/>
<line class="edge" x1="665" y1="166" x2="700" y2="228"/>
</svg>

Dependencies point inward: the domain layer depends on nothing, and gateways implement interfaces defined by the services that use them.

## Result

- Clear service boundaries that make it obvious where new logic belongs.
- Better maintainability and faster feature delivery across lending products.
- Shared patterns reinforced through code reviews and design discussions with the engineers I mentor.
```

- [ ] **Step 7: `content/case-studies/nbfc-co-lending.md`**

```markdown
---
title: Co-Lending Integrations with 4 NBFC Partners
summary: End-to-end integrations with Tata Capital, OFL, ABCL and IFMR that grew co-lending loan volume by 20%.
order: 2
tags: [go, integrations, fintech]
metrics: ["+20% loan volume", "4 NBFC partners", "Fund-flow visibility"]
---
## Problem

Co-lending lets a platform originate loans together with NBFC partners. Every new partner widens the lending distribution network, but each one brings its own integration, eligibility rules and fund-flow behaviour, and all of it has to stay RBI-compliant. When a loan wasn't assigned to a partner, working out why was slow.

## Approach

- Owned the integrations of **Tata Capital, OFL, ABCL and IFMR** end to end, from partner onboarding through the full loan lifecycle.
- Kept partner-specific behaviour behind a common partner contract, so the loan journey talks to one interface regardless of the lender.
- Added **fund assignment events** that record fund rejections along with the specific rule that failed, which turned fund-flow debugging into a lookup instead of an investigation.
- Worked closely with product, finance, compliance and operations to ship each integration on schedule.

## Architecture

<svg class="diagram" viewBox="0 0 770 300" role="img" aria-label="The loan journey calls fund allocation, which emits fund events and routes to a partner gateway with adapters for four NBFC partners">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="20" y="16" width="170" height="56" rx="10"/><text x="105" y="44">Loan journey</text>
<rect class="box" x="230" y="16" width="200" height="56" rx="10"/><text x="330" y="36">Fund allocation</text><text class="label" x="330" y="56">rules + fund events</text>
<rect class="box alt" x="470" y="16" width="200" height="56" rx="10"/><text x="570" y="36">Partner gateway</text><text class="label" x="570" y="56">common contract</text>
<rect class="box alt" x="20" y="230" width="160" height="48" rx="10"/><text x="100" y="254">Tata Capital</text>
<rect class="box alt" x="210" y="230" width="160" height="48" rx="10"/><text x="290" y="254">OFL</text>
<rect class="box alt" x="400" y="230" width="160" height="48" rx="10"/><text x="480" y="254">ABCL</text>
<rect class="box alt" x="590" y="230" width="160" height="48" rx="10"/><text x="670" y="254">IFMR</text>
<line class="edge" x1="190" y1="44" x2="226" y2="44"/>
<line class="edge" x1="430" y1="44" x2="466" y2="44"/>
<line class="edge" x1="540" y1="72" x2="110" y2="226"/>
<line class="edge" x1="555" y1="72" x2="295" y2="226"/>
<line class="edge" x1="575" y1="72" x2="482" y2="226"/>
<line class="edge" x1="595" y1="72" x2="665" y2="226"/>
</svg>

## Result

- **+20% co-lending loan volume** from the expanded partner network.
- One integration pattern reused across partners.
- Faster root-cause analysis for fund rejections.
```

- [ ] **Step 8: `content/case-studies/loan-documents.md`**

```markdown
---
title: Loan Document Generation Service
summary: A scalable, low-latency service that renders RBI-compliant loan agreements in 10 languages from HTML templates through headless Chrome.
order: 3
tags: [go, pdf, compliance]
metrics: ["10 languages", "RBI-compliant", "Low latency"]
---
## Problem

Every disbursed loan needs a legally compliant agreement, generated in the middle of the loan journey and in the borrower's language. Documents vary by product and partner, have to be exactly right for compliance, and must be fast enough not to slow the journey down.

## Approach

I designed and built a dedicated document generation service:

- **HTML templates** per document and language, filled with loan data in Go.
- **Headless Chrome** turns the rendered HTML into PDF. A browser engine gives correct text shaping for Indian scripts, so all 10 languages render properly without a custom typesetting stack.
- Built as its own service so generation scales independently of the rest of the lending stack.

## Architecture

<svg class="diagram" viewBox="0 0 770 130" role="img" aria-label="Loan data and a language template are rendered to HTML in Go, converted by headless Chrome, and produce a PDF agreement">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="10" y="36" width="125" height="56" rx="10"/><text x="72" y="64">Loan data</text>
<rect class="box" x="165" y="36" width="125" height="56" rx="10"/><text x="227" y="56">Template</text><text class="label" x="227" y="76">10 languages</text>
<rect class="box" x="320" y="36" width="125" height="56" rx="10"/><text x="382" y="56">Go renderer</text><text class="label" x="382" y="76">HTML</text>
<rect class="box alt" x="475" y="36" width="125" height="56" rx="10"/><text x="537" y="56">Headless</text><text x="537" y="76">Chrome</text>
<rect class="box alt" x="630" y="36" width="130" height="56" rx="10"/><text x="695" y="56">Agreement</text><text class="label" x="695" y="76">PDF</text>
<line class="edge" x1="135" y1="64" x2="161" y2="64"/>
<line class="edge" x1="290" y1="64" x2="316" y2="64"/>
<line class="edge" x1="445" y1="64" x2="471" y2="64"/>
<line class="edge" x1="600" y1="64" x2="626" y2="64"/>
</svg>

## Result

- RBI-compliant loan agreements in **10 languages**.
- A scalable, low-latency pipeline that fits inside the loan journey.
- Template-driven, so new documents and languages are content changes rather than new code paths.
```

- [ ] **Step 9: `assets/og-card.html`** and render it to `static/img/og.png`

```html
<!doctype html>
<html><head><meta charset="utf-8"><style>
@font-face { font-family: Inter; src: url("../static/fonts/inter-var.woff2"); font-weight: 100 900; }
html, body { margin: 0; width: 1200px; height: 630px; overflow: hidden; background: #0b1020; font-family: Inter, system-ui, sans-serif; color: #e6e8f2; }
.blob { position: absolute; border-radius: 50%; filter: blur(90px); opacity: .55; }
.card { position: absolute; inset: 70px; border-radius: 32px; background: rgba(255,255,255,.07); border: 1px solid rgba(255,255,255,.14); display: flex; flex-direction: column; justify-content: center; padding: 0 80px; }
h1 { font-size: 104px; margin: 0; letter-spacing: -.03em; }
p { font-size: 36px; margin: 12px 0 0; color: #b9c0d9; }
.tags { margin-top: 36px; font: 600 26px ui-monospace, Menlo, monospace; color: #5eead4; }
</style></head><body>
<div class="blob" style="width:620px;height:620px;background:#6366f1;left:-140px;top:-220px"></div>
<div class="blob" style="width:560px;height:560px;background:#8b5cf6;right:-120px;top:40px"></div>
<div class="blob" style="width:520px;height:520px;background:#14b8a6;left:300px;bottom:-300px"></div>
<div class="card"><h1>Prashant</h1><p>Senior Software Engineer · Backend Systems</p><div class="tags">Go · Distributed Systems · DDD · AWS</div></div>
</body></html>
```

```bash
"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --headless=new --disable-gpu --hide-scrollbars \
  --window-size=1200,630 --screenshot="$PWD/static/img/og.png" "file://$PWD/assets/og-card.html"
sips -g pixelWidth -g pixelHeight static/img/og.png   # expect 1200 x 630
```

- [ ] **Step 10: Build** `go run ./cmd/build` → `built dist in …`. Fix any validation output.
- [ ] **Step 11: Commit** — `git add content assets static/img/og.png && git commit -m "content: add Prashant's profile, experience and case studies"`

---

### Task 12: README, DEPLOY.md, workflow

**Files:** Create `README.md`, `DEPLOY.md`, `.github/workflows/deploy.yml`

- [ ] **Step 1: `.github/workflows/deploy.yml`**

```yaml
name: Deploy
on:
  push:
    branches: [main]
  pull_request:
  workflow_dispatch:

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: pages
  cancel-in-progress: false

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      - run: go vet ./...
      - run: go test ./...
      - run: go run ./cmd/build
      - if: github.event_name != 'pull_request'
        uses: actions/upload-pages-artifact@v3
        with:
          path: dist
  deploy:
    if: github.event_name != 'pull_request'
    needs: build
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - id: deployment
        uses: actions/deploy-pages@v4
```

- [ ] **Step 2: `README.md`**

````markdown
# prashant263.github.io

My portfolio, generated by a small static site generator written in Go.

```bash
go run ./cmd/serve    # http://localhost:8080, rebuilds on every save
go run ./cmd/build    # writes dist/
go test ./...
```

## Editing content

Everything lives in `content/`. No code changes needed:

| To… | Edit |
|---|---|
| Add a blog link | `content/blog.yaml` (example inside) |
| Update the resume | replace `static/resume.pdf`, bump `resume_updated` in `content/profile.yaml` |
| Change the photo | add `static/img/profile.jpg`, set `photo: img/profile.jpg` |
| Add a case study | new `content/case-studies/<slug>.md` with front matter |
| Add a project | `content/projects.yaml` |

The build stops on any content mistake and names the file and field.

## Layout

- `internal/content`: domain types, loading, validation (no HTML)
- `internal/render`: templates and Markdown to HTML (no disk writes)
- `internal/site`: builds `dist/`, copies static files, checks links
- `cmd/build`, `cmd/serve`: entry points
````

- [ ] **Step 3: `DEPLOY.md`**

````markdown
# Deploying to GitHub Pages

Do this on your personal laptop. Total cost: ₹0.

## 1. Prerequisites
- Go 1.24+ (`go version`) and Git (`git --version`).
- Logged in to GitHub as **prashant263**.

## 2. Copy the project
Copy this folder to your laptop (zip, AirDrop or a USB drive). Skip `dist/`.
For a clean history under your own Git identity:
```bash
cd portfolio
rm -rf .git
git init -b main
git add -A
git commit -m "Initial portfolio"
```

## 3. Personalise before publishing
In `content/profile.yaml`:
- **LinkedIn:** set the `url` of the LinkedIn social to your profile URL.
- **Photo:** save a square photo (≥ 800×800) as `static/img/profile.jpg` and set `photo: img/profile.jpg`.
- **Contact form:** sign up free at https://web3forms.com with your email, and paste the access key into `web3forms_key`.
- **Resume privacy:** `static/resume.pdf` is public and currently shows your phone number. Replace it with a version without the number if you prefer.

Check locally: `go run ./cmd/serve` and open http://localhost:8080.

## 4. Create the repository
1. On GitHub: **New repository**.
2. The name must be exactly **`prashant263.github.io`**.
3. Make it **Public**, and don't add a README, .gitignore or license.

## 5. Push
```bash
git remote add origin https://github.com/prashant263/prashant263.github.io.git
git push -u origin main
```

## 6. Turn on Pages
Repo → **Settings → Pages → Build and deployment → Source: GitHub Actions**.

Then open the **Actions** tab. The "Deploy" run should go green in about a minute. If it ran before you changed the setting and failed, click **Re-run all jobs**.

Your site is live at **https://prashant263.github.io**.

## Everyday updates
Edit files in `content/` or `static/`, then:
```bash
go run ./cmd/serve          # preview
git add -A && git commit -m "Add blog link" && git push
```
The site updates about a minute after each push.

## Optional: custom domain
Buy a domain, add `static/CNAME` containing just the domain (e.g. `prashant.dev`), set `site_url` in `profile.yaml`, push, then follow GitHub's "Managing a custom domain" guide for the DNS records.

## Troubleshooting
- **Build failed in Actions:** open the failed step. Content errors name the file and field, e.g. `blog.yaml: [0].date: is required`.
- **404 at the URL:** check the repo name is exactly `prashant263.github.io` and Pages source is "GitHub Actions".
- **Action versions deprecated:** bump the `@vN` numbers in `.github/workflows/deploy.yml` to the latest shown on each action's GitHub page.
````

- [ ] **Step 4: Commit** — `git add README.md DEPLOY.md .github && git commit -m "docs: add README, deploy guide and Pages workflow"`

---

### Task 13: Final verification

- [ ] **Step 1:** `gofmt -l . && go vet ./... && go test ./...` → no gofmt output, all tests pass.
- [ ] **Step 2:** Start the server in the background: `go run ./cmd/serve`. Wait for `serving dist at http://localhost:8080`.
- [ ] **Step 3: Smoke-test routes**

```bash
for p in / /blog/ /case-studies/ddd-migration/ /case-studies/nbfc-co-lending/ /case-studies/loan-documents/ /resume.pdf /sitemap.xml /nope; do
  printf "%s %s\n" "$(curl -s -o /dev/null -w '%{http_code}' "http://localhost:8080$p")" "$p"
done
```
Expected: 200 for all except `/nope` → 404.

- [ ] **Step 4: Visual check** (desktop + mobile screenshots, then Read the PNGs)

```bash
CH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
"$CH" --headless=new --disable-gpu --hide-scrollbars --window-size=1440,4000 --screenshot=/tmp/desktop.png http://localhost:8080/
"$CH" --headless=new --disable-gpu --hide-scrollbars --window-size=390,6000 --screenshot=/tmp/mobile.png http://localhost:8080/
"$CH" --headless=new --disable-gpu --hide-scrollbars --window-size=1440,1800 --screenshot=/tmp/case.png http://localhost:8080/case-studies/ddd-migration/
```
Check: mesh background visible, glass cards, photo ring, no horizontal overflow on mobile, diagram legible. The tall window puts every section in the viewport, so `.reveal` sections should be visible; if any look faded, add `--virtual-time-budget=3000` so the transitions finish before capture.

- [ ] **Step 5: Rebuild-on-save check:** edit a word in `content/profile.yaml`, wait 1s, `curl -s localhost:8080/ | grep` for it, then revert.
- [ ] **Step 6:** Stop the server. Final commit if anything changed.
