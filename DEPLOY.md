# Make the portfolio live — step by step

Goal: your site running at **https://prashant263.github.io**, deployed from your personal laptop.
Cost: **₹0** (GitHub Pages + GitHub Actions on a public repo + Google Forms).
Time: about 30 minutes the first time.

```
This work laptop ──zip──▶ Personal laptop ──git push──▶ GitHub ──Actions builds──▶ https://prashant263.github.io
```

Every later update is just: edit a file → `git push` → live in about a minute.

---

## Part A — On this (work) laptop: package the folder

**A1. Zip the project** (leaves out the build output, local settings and this machine's git history):

```bash
cd ~/Desktop
zip -r portfolio.zip temp -x "temp/dist/*" "temp/.claude/*" "temp/.git/*" "*.DS_Store"
ls -lh portfolio.zip          # about 1 MB
```

**A2. Move `portfolio.zip` to your personal laptop**: Google Drive, email to yourself, AirDrop, or a USB drive.

That's everything needed from this machine. Nothing here touches GitHub.

---

## Part B — On your personal laptop: one-time setup

**B1. Install the tools** (skip any you already have):

| Tool | macOS | Windows |
|---|---|---|
| Go **1.24.5 or newer** | `brew install go` or https://go.dev/dl | installer from https://go.dev/dl |
| Git | `xcode-select --install` or `brew install git` | https://git-scm.com/download/win |
| GitHub CLI (recommended) | `brew install gh` | `winget install GitHub.cli` |

Check they work (open a new terminal first):

```bash
go version        # go1.24.5 or higher
git --version
gh --version      # if you installed it
```

**B2. Tell Git who you are** (used on your commits):

```bash
git config --global user.name  "Prashant B"
git config --global user.email "prashantab8505@gmail.com"
```

**B3. Unzip and rename the folder** so it matches the repo name:

```bash
cd ~/Desktop                         # wherever you saved the zip
unzip portfolio.zip
mv temp prashant263.github.io
cd prashant263.github.io
```

**B4. Check it builds and runs locally:**

```bash
go test ./...                        # every line should say ok
go run ./cmd/serve                   # leave this running
```

Open **http://localhost:8080** in your browser and click around. Stop the server with `Ctrl+C` when you're done.
(If port 8080 is taken, use `go run ./cmd/serve -addr localhost:8090` and open that port instead.)

---

## Part C — Final check before it's public

Once pushed, anyone can see these. Go through the list:

- [ ] **Phone number.** `static/resume.pdf` and its preview image `static/img/resume-preview.png` both show your phone number. To hide it: export a resume PDF without the number, save it over `static/resume.pdf`, then regenerate the preview:
  - macOS: `sips -s format png --resampleWidth 600 static/resume.pdf --out static/img/resume-preview.png`
  - Windows: open the PDF, screenshot page 1, save it as `static/img/resume-preview.png` (about 600 px wide).
  - Then update `resume_updated:` in `content/profile.yaml`.
- [ ] **"Open to backend roles"** in the Contact section — edit or remove it in `templates/pages/home.html` if you don't want that public while employed.
- [ ] **Status pill** "Currently building at KreditBee" — `status:` in `content/profile.yaml` (delete the line to hide it).
- [ ] **Social links** under `socials:` in `content/profile.yaml` (GitHub, LinkedIn, Medium, email) all open the right profiles.
- [ ] **Contact form** — it's already connected to your "Connection Details" Google Form. Send one test message from http://localhost:8080/#contact and confirm it appears in the Form's **Responses** tab. (See [Contact form](#contact-form-google-form) if you ever need to reconnect it.)
- [ ] **Case studies** in `content/case-studies/` — only say what you're comfortable discussing in an interview.

After any change, run `go run ./cmd/serve` again to preview.

---

## Part D — Put it on GitHub

**D1. Start a fresh git history** (inside the `prashant263.github.io` folder):

```bash
git init -b main
git add -A
git commit -m "Initial portfolio"
```

**D2. Create the repository and push.** Pick one option.

**Option 1: GitHub CLI (easiest)**

```bash
gh auth login                        # choose GitHub.com → HTTPS → log in with a web browser
gh repo create prashant263.github.io --public --source=. --remote=origin --push
```

**Option 2: GitHub website**

1. Go to https://github.com/new while logged in as **prashant263**.
2. **Repository name:** exactly `prashant263.github.io` (it must match your username).
3. Choose **Public**. Do **not** add a README, .gitignore or license.
4. Click **Create repository**, then run:
   ```bash
   git remote add origin https://github.com/prashant263/prashant263.github.io.git
   git push -u origin main
   ```
5. When asked for a password, paste a **Personal Access Token**, not your GitHub password: GitHub → profile picture → **Settings → Developer settings → Personal access tokens → Tokens (classic) → Generate new token**, tick **repo** and **workflow**, copy the token.

---

## Part E — Switch on GitHub Pages

1. Open https://github.com/prashant263/prashant263.github.io → **Settings** → **Pages** (left sidebar).
2. Under **Build and deployment → Source**, choose **GitHub Actions**.
3. Open the **Actions** tab. You'll see a run called **Deploy**.
   - If it's red (it ran before Pages was enabled), click it → **Re-run all jobs**.
   - Wait for both jobs, **build** and **deploy**, to turn green (about 1–2 minutes).
4. Open **https://prashant263.github.io**. The very first deploy can take a few extra minutes to appear.

**🎉 Your site is live.**

---

## Part F — Check the live site

- [ ] Every section loads, the photo shows, the theme toggle works, and the P bubbles pop.
- [ ] **Download CV** downloads your resume.
- [ ] Nav links (About, Experience, Work, Resume, Blog, Contact) all jump to the right place.
- [ ] Send one real test message from the live Contact form, check it in Google Forms → **Responses**, then delete that test response.
- [ ] Phone check: open the site on your phone.
- [ ] Sharing preview: paste your URL into https://www.linkedin.com/post-inspector/ to see the card LinkedIn will show.

Then add the link everywhere: LinkedIn (Contact info → Website), your GitHub profile, and your resume.

---

## Everyday updates

Edit, preview, push:

```bash
go run ./cmd/serve                         # preview at http://localhost:8080 (Ctrl+C to stop)
git add -A
git commit -m "Add a blog link"
git push                                   # live about a minute later
```

| I want to… | Edit |
|---|---|
| Add a blog/Medium post | `content/blog.yaml` (copy the example inside) |
| Update my resume | replace `static/resume.pdf`, regenerate `static/img/resume-preview.png` (see Part C), bump `resume_updated` in `content/profile.yaml` |
| Change my photo | replace `static/img/profile.jpeg` (portrait 4:5 looks best) |
| Change the P bubbles | `bubbles:` in `content/profile.yaml` (keep each detail to about 60 characters) |
| Add a side project | `content/projects.yaml` (`live:` for a website, `repo:` for code) |
| Add an award or milestone | `content/recognition.yaml` |
| Change headline, pitch, About cards | `content/profile.yaml` |
| Add a case study | new file `content/case-studies/<name>.md` (copy an existing one) |

If you make a typo in a content file, the build stops and tells you the exact file and field, so a broken page can never go live. Fix it, then push again.

Note: if you change Go code (rare), stop `go run ./cmd/serve` and start it again. Content, template and CSS changes reload by themselves.

---

## Contact form (Google Form)

**Already set up** and connected to your "Connection Details" form. Every message becomes a row in that form's responses (and its Google Sheet, if linked). Use these steps only if you create a new form.

1. On the form's **Responses** tab, click **Link to Sheets** once to keep every message in a spreadsheet.
2. Turn on alerts: **Responses → ⋮ → Get email notifications for new responses**.
3. Keep **Settings → Responses** set so it does *not* require sign-in or collect email addresses. Otherwise the website can't submit to it.

To connect a different form:

1. Create the form with four questions in this order: **Name** (short answer, required), **Contact number** (short answer, optional), **Email** (short answer, required), **Message** (paragraph, required).
2. Click **⋮ → Get pre-filled link**, type `NAME`, `PHONE`, `EMAIL`, `MESSAGE` as the answers, click **Get link** and copy it. It contains `.../forms/d/e/<form-id>/viewform?...&entry.111=NAME&entry.222=PHONE&entry.333=EMAIL&entry.444=MESSAGE`.
3. Update `contact_form:` in `content/profile.yaml`:
   ```yaml
   contact_form:
     action: "https://docs.google.com/forms/d/e/<form-id>/formResponse"   # same ID, ends in /formResponse
     name: "entry.111"
     phone: "entry.222"
     email: "entry.333"
     message: "entry.444"
   ```
4. Preview, send a test, check that it shows up in Responses, then push.

The website can't read Google's reply, so checking the Responses tab is the only way to confirm a message arrived.

---

## Optional: your own domain (about ₹800–1000 a year)

1. Buy a domain (e.g. `prashant.dev`) from any registrar.
2. Repo → **Settings → Pages → Custom domain**, enter it, and click **Save**.
3. At your registrar, add the DNS records shown in GitHub's guide: https://docs.github.com/pages/configuring-a-custom-domain-for-your-github-pages-site
4. Set `site_url:` in `content/profile.yaml` to `https://your-domain`, then commit and push.
5. Back in **Settings → Pages**, tick **Enforce HTTPS** once it becomes available. This can take up to a day.

---

## Troubleshooting

| Problem | Fix |
|---|---|
| `git push` asks for a password and rejects it | Use the GitHub CLI (`gh auth login`) or a Personal Access Token (Part D, Option 2 step 5). Your normal password doesn't work for git. |
| Actions run is red at **build** | Click the run → the red step. Content mistakes name the file and field, e.g. `blog.yaml: [0].date: is required`. Fix it, commit, push. |
| Actions run is red at **deploy** | Check **Settings → Pages → Source** is **GitHub Actions**, then **Re-run all jobs**. |
| No Actions run appeared after pushing | Check `.github/workflows/deploy.yml` exists in the repo and the branch is `main`. Repo → **Actions** → enable workflows if GitHub asks. |
| Site shows 404 | Repo name must be exactly `prashant263.github.io` and public. Wait 5–10 minutes after the first green deploy. |
| Old version still showing | Hard-refresh (`Cmd+Shift+R` / `Ctrl+F5`). Updates usually appear within a minute or two. |
| `go: go.mod requires go >= 1.24.5` | Update Go from https://go.dev/dl. |
| A GitHub warning that an action is deprecated | Bump the `@vN` numbers in `.github/workflows/deploy.yml` to the latest version listed on each action's GitHub page. |
| Contact form says sent but nothing appears in Responses | Form must not require sign-in. Field IDs in `contact_form:` must match the form (redo the pre-filled link step). |
