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
- **LinkedIn:** set the `url` of the LinkedIn social to your profile URL.
- **Photo:** save a portrait photo (4:5, e.g. 800×1000 or larger) as `static/img/profile.jpg` and set `photo: img/profile.jpg`. It fills a rounded rectangle in the hero; other shapes are cropped from the centre.
- **Contact form:** sign up free at https://web3forms.com with your email, and paste the access key into `web3forms_key`.
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

## Optional: custom domain
Buy a domain, set it in repo **Settings → Pages → Custom domain** (with Actions-based Pages a `CNAME` file isn't used), update `site_url` in `profile.yaml`, push, then follow GitHub's "Managing a custom domain" guide for the DNS records.

## Troubleshooting
- **Build failed in Actions:** open the failed step. Content errors name the file and field, e.g. `blog.yaml: [0].date: is required`.
- **404 at the URL:** check the repo name is exactly `prashant263.github.io` and Pages source is "GitHub Actions".
- **Action versions deprecated:** bump the `@vN` numbers in `.github/workflows/deploy.yml` to the latest shown on each action's GitHub page.
