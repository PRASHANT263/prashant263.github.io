# Deploying to GitHub Pages

Do this on your personal laptop. Total cost: ₹0.

## 1. Prerequisites
- Go 1.24.5+ (`go version`) and Git (`git --version`).
- Logged in to GitHub as **prashant263**.

## 2. Copy the project
First make sure the finished work is on `main`. In the original folder, run `git checkout main && git merge feat/portfolio` (skip this if it's already merged).

Copy the folder to your laptop (zip, AirDrop or a USB drive). Skip `dist/`, `.claude/` and `.screenshots/`; they're local-only and gitignored anyway.
For a clean history under your own Git identity:
```bash
cd <the-copied-folder>
rm -rf .git
git init -b main
git add -A
git commit -m "Initial portfolio"
```

## 3. Personalise before publishing
In `content/profile.yaml`:
- **Social links:** GitHub, LinkedIn and Medium are set under `socials:`; double-check the URLs (icon can be github, linkedin, medium or email).
- **Photo:** save a portrait photo (4:5, e.g. 800×1000 or larger) as `static/img/profile.jpg` and set `photo: img/profile.jpg`. It fills a rounded rectangle in the hero; other shapes are cropped from the centre.
- **Contact form:** messages go to a Google Form and its Google Sheet. Set up `contact_form:` in `profile.yaml` with the steps in [Contact form (Google Form)](#contact-form-google-form) below.
- **Resume privacy:** `static/resume.pdf` **and** its preview image `static/img/resume-preview.png` are public, and both currently show your phone number. If you'd rather not publish it, replace the PDF with a version without the number, then regenerate the preview:
  `sips -s format png --resampleWidth 600 static/resume.pdf --out static/img/resume-preview.png`
- **Availability line:** the Contact section says "Open to backend roles". Edit or remove it in `templates/pages/home.html` if you don't want that public.
- **Case studies:** `content/case-studies/*.md` stick to what's on your resume. Add detail you're comfortable discussing in interviews.

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

## Contact form (Google Form)
Messages from the site's contact form are saved as rows in a Google Sheet. Free, no account keys.

1. Go to https://forms.google.com and create a blank form, e.g. "Portfolio contact".
2. Add four questions, in this order:
   - **Name**: Short answer, Required
   - **Contact number**: Short answer (optional)
   - **Email**: Short answer, Required
   - **Message**: Paragraph, Required
3. **Settings** tab: under *Responses*, make sure **Collect email addresses** is *Do not collect*, and **Restrict to users in your organisation** / **Limit to 1 response** are *off*. (If sign-in is required, the site can't submit.)
4. **Responses** tab: click **Link to Sheets** and create a new spreadsheet. Every message lands there.
5. Find the field IDs: click **⋮ → Get pre-filled link**, type `NAME`, `PHONE`, `EMAIL`, `MESSAGE` into the four answers, click **Get link** and copy it. It looks like:
   `https://docs.google.com/forms/d/e/1FAIpQLSe.../viewform?usp=pp_url&entry.123=NAME&entry.456=PHONE&entry.789=EMAIL&entry.101=MESSAGE`
6. Fill in `contact_form:` in `content/profile.yaml`:
   ```yaml
   contact_form:
     action: "https://docs.google.com/forms/d/e/1FAIpQLSe.../formResponse"   # same ID, but end with /formResponse
     name: "entry.123"
     phone: "entry.456"
     email: "entry.789"
     message: "entry.101"
   ```
7. Test it: run `go run ./cmd/serve`, send a message from http://localhost:8080/#contact, and check that a new row appears in the Sheet within a few seconds. (The site can't read Google's reply, so this check is the way to confirm it works.)

Tip: in the Form, **Responses → ⋮ → Get email notifications for new responses** emails you whenever someone writes.

## Optional: custom domain
Buy a domain, set it in repo **Settings → Pages → Custom domain** (with Actions-based Pages a `CNAME` file isn't used), update `site_url` in `profile.yaml`, push, then follow GitHub's "Managing a custom domain" guide for the DNS records.

## Troubleshooting
- **Build failed in Actions:** open the failed step. Content errors name the file and field, e.g. `blog.yaml: [0].date: is required`.
- **404 at the URL:** check the repo name is exactly `prashant263.github.io` and Pages source is "GitHub Actions".
- **Action versions deprecated:** bump the `@vN` numbers in `.github/workflows/deploy.yml` to the latest shown on each action's GitHub page.
