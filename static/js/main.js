(() => {
  const root = document.documentElement;

  // Theme: an explicit choice overrides the OS preference and is remembered.
  const themeBtn = document.querySelector(".theme-toggle");
  if (themeBtn) {
    themeBtn.addEventListener("click", () => {
      const stored = root.dataset.theme;
      const current = (stored === "light" || stored === "dark") ? stored : (matchMedia("(prefers-color-scheme: light)").matches ? "light" : "dark");
      const next = current === "light" ? "dark" : "light";
      root.dataset.theme = next;
      try { localStorage.setItem("theme", next); } catch (e) { /* storage unavailable */ }
    });
  }

  // Mobile nav.
  const navToggle = document.querySelector(".nav-toggle");
  const navLinks = document.getElementById("nav-links");
  if (navToggle && navLinks) {
    const setOpen = (open) => {
      navLinks.classList.toggle("open", open);
      navToggle.setAttribute("aria-expanded", String(open));
    };
    navToggle.addEventListener("click", () => setOpen(!navLinks.classList.contains("open")));
    navLinks.addEventListener("click", (e) => { if (e.target.closest("a")) setOpen(false); });
  }

  // Hero greeting: type and delete each greeting in turn (Hi → नमस्ते → ನಮಸ್ಕಾರ).
  // Splits by grapheme so Devanagari and Kannada letters never break apart.
  const greet = document.querySelector(".greet[data-greetings]");
  const words = greet ? greet.dataset.greetings.split("|").filter(Boolean) : [];
  if (words.length > 1 && !matchMedia("(prefers-reduced-motion: reduce)").matches) {
    const graphemes = (s) => ("Segmenter" in Intl)
      ? Array.from(new Intl.Segmenter().segment(s), (g) => g.segment)
      : Array.from(s);
    const letters = words.map(graphemes);
    let word = 0, shown = letters[0].length, deleting = true;
    const tick = () => {
      shown += deleting ? -1 : 1;
      greet.textContent = letters[word].slice(0, shown).join("");
      let delay = deleting ? 70 : 140;
      if (deleting && shown === 0) {
        deleting = false;
        word = (word + 1) % letters.length;
        delay = 350;
      } else if (!deleting && shown === letters[word].length) {
        deleting = true;
        delay = 1800;
      }
      setTimeout(tick, delay);
    };
    setTimeout(tick, 1800);
  }

  // Reveal sections as they scroll into view.
  root.classList.add("reveal-ready");
  const reveals = document.querySelectorAll(".reveal");
  if ("IntersectionObserver" in window) {
    const io = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        if (entry.isIntersecting) {
          entry.target.classList.add("visible");
          io.unobserve(entry.target);
        }
      }
    }, { threshold: 0.08 });
    reveals.forEach((el) => io.observe(el));
  } else {
    reveals.forEach((el) => el.classList.add("visible"));
  }

  // Blog tag filter.
  const filter = document.querySelector(".tag-filter");
  const list = document.querySelector("[data-blog-list]");
  if (filter && list) {
    filter.addEventListener("click", (e) => {
      const btn = e.target.closest("button[data-tag]");
      if (!btn) return;
      const tag = btn.dataset.tag;
      filter.querySelectorAll("button").forEach((b) => {
        b.classList.toggle("active", b === btn);
        b.setAttribute("aria-pressed", String(b === btn));
      });
      list.querySelectorAll("[data-tags]").forEach((card) => {
        card.hidden = tag !== "" && !card.dataset.tags.split(" ").includes(tag);
      });
    });
  }

  // Contact form via Web3Forms; without JS the form posts normally.
  const form = document.querySelector("form[data-contact]");
  if (form) {
    const status = form.querySelector(".form-status");
    const button = form.querySelector("button[type=submit]");
    const show = (msg, kind) => { status.textContent = msg; status.dataset.kind = kind; };
    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!form.reportValidity()) return;
      const data = Object.fromEntries(new FormData(form));
      if (data.botcheck) {
        show("Thanks! Your message has been sent.", "success");
        return;
      }
      delete data.redirect;
      button.disabled = true;
      show("Sending…", "info");
      try {
        const res = await fetch(form.action, {
          method: "POST",
          headers: { "Content-Type": "application/json", Accept: "application/json" },
          body: JSON.stringify(data),
          signal: AbortSignal.timeout ? AbortSignal.timeout(15000) : undefined,
        });
        const body = await res.json().catch(() => ({}));
        if (!res.ok || body.success === false) throw new Error(body.message || "request failed");
        form.reset();
        show("Thanks! Your message has been sent.", "success");
      } catch (err) {
        show("Sorry, the message couldn't be sent. Please email me directly.", "error");
      } finally {
        button.disabled = false;
      }
    });
  }
})();
