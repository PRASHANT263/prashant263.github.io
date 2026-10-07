# Portfolio Website — Design Spec

**Date:** 2026-10-07
**Owner:** Prashant (Senior Software Engineer, Backend Systems)
**Live URL:** https://prashantb.github.io
**Cost:** ₹0 (GitHub Pages + Web3Forms free tier + GitHub Actions on a public repo)

## 1. Goals

One site serving three purposes:

1. **Job hunting** — recruiters/hiring managers quickly see impact, case studies, and can download the resume.
2. **Personal** — a long-term home for blog links, projects, and contact.
3. **Engineering showcase** — the site's own source (a custom Go static site generator with layered design, tests, CI) is itself a portfolio piece.

Non-goals: CMS/admin UI, database, server runtime, analytics, comments, i18n.

## 2. Approach

A custom **Go static site generator** reads content files (YAML + Markdown), renders HTML with `html/template` and `goldmark`, and writes a `dist/` folder. GitHub Actions builds and deploys `dist/` to GitHub Pages on every push to `main`. The contact form posts directly to **Web3Forms** (no backend).

All content edits (adding a blog link, updating the resume) are file edits + `git push`; no Go code changes required.

## 3. Repository Structure

Repo name: `prashantb.github.io` (GitHub user `prashantb`).

```
.
├── go.mod                    # module github.com/prashantb/prashantb.github.io
├── content/
│   ├── profile.yaml          # name, title, pitch, about, quick facts, socials, photo, resume, web3forms key
│   ├── experience.yaml       # company → roles → bullets + highlight chips
│   ├── skills.yaml           # grouped skill tags
│   ├── projects.yaml         # side projects
│   ├── blog.yaml             # blog entries (external links; local posts supported)
│   ├── case-studies/*.md     # front matter + Markdown body
│   └── posts/*.md            # optional local blog posts (empty at launch)
├── internal/
│   ├── content/              # domain types, loaders, validation (no HTML, no templates)
│   ├── render/               # domain → HTML via templates + goldmark (no filesystem writes)
│   └── site/                 # orchestration: load → validate → render → write dist/, copy static
├── templates/                # base layout, partials, page templates
├── static/
│   ├── css/  js/  fonts/
│   ├── img/profile.jpg       # profile photo (provided by owner)
│   └── resume.pdf            # downloadable resume
├── cmd/
│   ├── build/                # `go run ./cmd/build` → dist/
│   └── serve/                # local dev server on :8080, rebuilds on file change
└── .github/workflows/deploy.yml
```

The existing empty `main.go` at the repo root is removed; entry points live under `cmd/`.

**Dependency rule:** `site` → `render` → `content`. `content` knows nothing about HTML; `render` never touches disk; only `site` and `cmd/*` do I/O.

## 4. Content Model

### `blog.yaml`
```yaml
- title: "Migrating a lending monolith to DDD"
  date: 2026-09-12          # YYYY-MM-DD, required
  url: https://medium.com/@prashant/...   # required for external entries
  source: Medium            # free text badge: Medium | Dev.to | Hashnode | LinkedIn | ...
  tags: [go, ddd]
  summary: "Lessons from splitting domain, service and gateway layers."
```
An entry is **external** when `url` is set, **local** when `slug` is set (renders `content/posts/<slug>.md` at `/blog/<slug>/`). Exactly one of `url`/`slug` must be present.

### Case studies (`content/case-studies/<slug>.md`)
Front matter: `title`, `summary`, `order`, `tags`, `metrics` (list of short chip strings). Body sections: Problem → Approach → Architecture (inline SVG diagram) → Result.

Launch case studies (sanitized — only figures already on the public resume, no internal service names, partner API details, or confidential data):
1. Monolith → DDD migration of lending systems
2. NBFC co-lending integrations (4 partners, +20% loan volume)
3. Loan document generation service (HTML templates + headless Chrome, 10 languages)

### Validation (build fails with file + field in the message)
- Required fields present; dates parse as `YYYY-MM-DD`; URLs are absolute `http(s)`.
- Blog: exactly one of `url`/`slug`; local slug file exists.
- Referenced files exist: profile photo, resume PDF, case-study/post Markdown.
- No duplicate slugs.

## 5. Pages & Sections

Single-scroll home page (`/`), in order:

1. **Hero** — circular profile photo with animated gradient ring; name; "Senior Software Engineer · Backend Systems"; one-line pitch; buttons *View work* and *Download resume*; GitHub / LinkedIn / email icons.
2. **About** — short narrative + quick facts (4+ years, 3× Star Performer, Bengaluru).
3. **Experience** — vertical timeline: Intern → SWE → Senior SWE at KreditBee, bullets plus highlight chips (+20% loan volume, 10 languages, −30% toil, mentoring 5 engineers).
4. **Case Studies** — three glass cards linking to `/case-studies/<slug>/`.
5. **Skills** — grouped tag clouds from `skills.yaml`.
6. **Resume** — glass card with page-1 preview image, *Download PDF* (`download` attribute) and *View online* (new tab) buttons, and "Last updated" date from `profile.yaml`.
7. **Blog** — latest 6 entries (date desc): source badge, date, tags, summary; external links open in new tab with `rel="noopener"`. *See all* → `/blog/`.
8. **Projects** — cards from `projects.yaml` with GitHub/live links.
9. **Contact** — glass form (name, email, message) + direct `mailto:` link.

Other pages: `/blog/` (all entries, client-side tag filter), `/case-studies/<slug>/`, `/blog/<slug>/` (local posts only), `404.html`.

Sticky glass nav with anchor links to sections; collapses to a menu on mobile.

## 6. Visual Design

- **Background:** gradient mesh — 3–4 large blurred blobs (indigo `#6366f1`, violet `#8b5cf6`, teal `#14b8a6`) drifting slowly via CSS keyframes. No JS. Animation disabled under `prefers-reduced-motion`.
- **Surfaces:** glass cards — semi-transparent fill, `backdrop-filter: blur(16px)`, 1px translucent border, soft shadow, 16px radius.
- **Themes:** dark (default, deep navy base `#0b1020`) and light (soft off-white base, pastel blobs). Toggle in nav; respects `prefers-color-scheme` on first visit; choice saved in `localStorage` (wrapped in try/catch).
- **Type:** Inter (body/headings), JetBrains Mono (tags, chips, small accents), self-hosted WOFF2.
- **Motion:** subtle fade/slide reveal on scroll via `IntersectionObserver`; none under reduced motion.
- **Responsive:** mobile-first; single column under 768px; no horizontal scroll.
- **JS budget:** < 10KB total, vanilla: theme toggle, scroll reveal, mobile nav, blog tag filter, contact form submit.

## 7. Contact Form

- Form `POST`s via `fetch` to `https://api.web3forms.com/submit` with the public access key from `profile.yaml`.
- Fields: name (required, ≤100), email (required, `type=email`), message (required, 10–5000 chars), hidden honeypot `botcheck`.
- Client-side validation before submit; inline success/error message; inputs preserved on failure; button disabled while sending.
- Without JS, the form submits normally to Web3Forms with a `redirect` back to `/#contact`. The `mailto:` link is always visible as a fallback.
- Free tier: 250 submissions/month.

## 8. SEO & Sharing

Per-page `<title>` and meta description; Open Graph + Twitter card tags with a static 1200×630 OG image (`static/img/og.png`, built from the hero design); `sitemap.xml` and `robots.txt` generated at build; canonical URLs under `https://prashantb.github.io`.

## 9. Error Handling

| Layer | Behavior |
|---|---|
| Build | Any load/validation/template error aborts with a non-zero exit and a message naming file + field. Broken content never deploys. |
| Internal links | Post-build checker verifies every internal `href`/`src` in `dist/` resolves; failure aborts deploy. |
| Contact form | Network/API failure shows a visible error and the `mailto:` fallback; never fails silently. |

## 10. Testing

- `internal/content`: table-driven tests for parsing and every validation rule (valid + each invalid case).
- `internal/render`: golden-file tests rendering fixture content and diffing against committed HTML snapshots (`-update` flag to refresh).
- `internal/site`: integration test building fixture content into a temp dir and asserting expected files exist.
- Link checker runs as part of `cmd/build` (or a `-check` flag) and in CI.
- CI: `go vet`, `go test ./...`, build, link check.

## 11. Deployment

GitHub Actions workflow `deploy.yml` on push to `main`:
1. Checkout, set up Go.
2. `go vet ./...` and `go test ./...`.
3. `go run ./cmd/build` → `dist/` (includes link check).
4. `actions/upload-pages-artifact` with `dist/`, then `actions/deploy-pages`.

Pull requests run steps 1–3 only. Repo Settings → Pages → Source: *GitHub Actions*.

## 12. Owner-Provided Assets

- `static/img/profile.jpg` — square profile photo, ≥ 800×800.
- `static/resume.pdf` — defaults to the existing `Prashants_a25.pdf` in the project root (note: the public file will expose whatever contact details it contains, including phone number; consider a version without the phone number).
- GitHub username confirmation (`prashantb`), LinkedIn URL, Web3Forms access key (free signup with email).
- Blog entries (optional at launch; section shows an empty-state message if none).
- Side projects (optional at launch; section hidden if `projects.yaml` is empty).

## 13. Future Extensions (not in scope)

Custom domain (CNAME file + DNS), local Markdown blog posts (already supported by the model), RSS feed, analytics.
