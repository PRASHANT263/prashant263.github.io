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
