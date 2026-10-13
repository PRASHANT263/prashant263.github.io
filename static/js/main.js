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

  // Skill bubbles: skills float up out of the avatar and fade away; clicking
  // one pops it with a few sparks and shows a short fact underneath.
  const avatar = document.querySelector(".hero-avatar[data-bubbles]");
  const motionOK = !matchMedia("(prefers-reduced-motion: reduce)").matches;
  if (avatar && motionOK) {
    let skills = [];
    try { skills = JSON.parse(avatar.querySelector(".bubble-data").textContent) || []; } catch (e) { /* bad data: no bubbles */ }
    const card = avatar.querySelector(".bubble-card");
    const rand = (min, max) => min + Math.random() * (max - min);
    const live = new Set();
    let next = 0, hideCard;

    const showFact = (skill) => {
      card.innerHTML = "";
      const title = document.createElement("strong");
      const text = document.createElement("p");
      title.textContent = skill.skill;
      text.textContent = skill.detail;
      card.append(title, text);
      card.classList.add("show");
      avatar.classList.add("used");
      clearTimeout(hideCard);
      hideCard = setTimeout(() => card.classList.remove("show"), 4500);
    };

    const sparks = (x, y) => {
      for (let i = 0; i < 10; i++) {
        const s = document.createElement("span");
        const angle = (i / 10) * Math.PI * 2;
        const dist = rand(26, 46);
        s.className = "spark";
        s.style.left = x + "px";
        s.style.top = y + "px";
        s.style.setProperty("--sx", Math.cos(angle) * dist + "px");
        s.style.setProperty("--sy", Math.sin(angle) * dist + "px");
        s.addEventListener("animationend", () => s.remove());
        avatar.append(s);
      }
    };

    const pop = (b, skill) => {
      const a = avatar.getBoundingClientRect(), r = b.getBoundingClientRect();
      b.style.transform = getComputedStyle(b).transform; // freeze where it is
      b.classList.add("popping");
      requestAnimationFrame(() => { b.style.transform += " scale(1.7)"; });
      setTimeout(() => { b.remove(); live.delete(b); }, 300);
      sparks(r.left - a.left + r.width / 2, r.top - a.top + r.height / 2);
      showFact(skill);
    };

    const spawn = () => {
      if (!skills.length || live.size >= 4 || document.hidden || !avatar.offsetParent) return;
      const skill = skills[next++ % skills.length];
      const b = document.createElement("button");
      b.type = "button";
      b.className = "bubble";
      b.textContent = skill.skill;
      b.setAttribute("aria-label", skill.skill + ": show detail");
      b.style.setProperty("--dx", rand(-110, 110) + "px");
      b.style.setProperty("--dy", rand(-230, -150) + "px");
      b.style.setProperty("--wob", rand(-18, 18) + "px");
      b.style.setProperty("--dur", rand(6, 8).toFixed(1) + "s");
      b.addEventListener("click", (e) => { e.stopPropagation(); pop(b, skill); });
      b.addEventListener("animationend", () => { b.remove(); live.delete(b); });
      live.add(b);
      avatar.append(b);
    };

    avatar.addEventListener("click", () => { for (let i = 0; i < 3; i++) setTimeout(spawn, i * 180); });
    setTimeout(spawn, 900);
    setInterval(spawn, 1800);
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

  // Contact form → Google Form (responses land in its Google Sheet).
  // Google doesn't let other sites read the reply (no CORS), so a request that
  // completes without a network error counts as sent. Without JS the form
  // posts normally and Google's confirmation opens in a new tab.
  const form = document.querySelector("form[data-contact]");
  if (form) {
    const status = form.querySelector(".form-status");
    const button = form.querySelector("button[type=submit]");
    const trap = form.querySelector("[data-honeypot]");
    const show = (msg, kind) => { status.textContent = msg; status.dataset.kind = kind; };
    const sent = () => { form.reset(); show("Thanks! Your message has been sent.", "success"); };
    form.addEventListener("submit", async (e) => {
      e.preventDefault();
      if (!form.reportValidity()) return;
      if (trap && trap.value) { sent(); return; } // only bots fill the hidden field
      if (form.hasAttribute("data-unconfigured")) {
        show("This form isn't connected yet. Please use the email link.", "error");
        return;
      }
      button.disabled = true;
      show("Sending…", "info");
      try {
        await fetch(form.action, {
          method: "POST",
          mode: "no-cors",
          body: new URLSearchParams(new FormData(form)),
          signal: AbortSignal.timeout ? AbortSignal.timeout(15000) : undefined,
        });
        sent();
      } catch (err) {
        show("Sorry, the message couldn't be sent. Please email me directly.", "error");
      } finally {
        button.disabled = false;
      }
    });
  }
})();
