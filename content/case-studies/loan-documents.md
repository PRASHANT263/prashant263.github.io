---
title: Loan Document Generation Service
summary: A scalable, low-latency service that generates RBI-compliant loan agreements in 10 languages using HTML templates and a headless Chrome PDF pipeline.
order: 3
tags: [go, pdf, compliance]
metrics: ["10 languages", "RBI-compliant", "Low latency"]
---
## Context

Every loan needs a compliant loan agreement, and borrowers need it in a language they can read. Documents are generated within the loan journey, so the service has to be fast and scalable as well as correct.

## What I did

I designed and built a dedicated Loan Document Generation service:

- **HTML templates** define each agreement, with versions for each of the 10 supported languages.
- A **headless Chrome PDF pipeline** turns the rendered HTML into the final PDF agreement.
- The service is built to be **scalable and low-latency**, and the agreements meet RBI requirements.

## Architecture

<svg class="diagram" viewBox="0 0 770 130" role="img" aria-label="Loan data fills an HTML template in one of ten languages, headless Chrome converts it, and the result is a PDF loan agreement">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="10" y="36" width="160" height="56" rx="10"/><text x="90" y="64">Loan data</text>
<rect class="box" x="210" y="36" width="160" height="56" rx="10"/><text x="290" y="56">HTML template</text><text class="label" x="290" y="76">10 languages</text>
<rect class="box alt" x="410" y="36" width="160" height="56" rx="10"/><text x="490" y="56">Headless</text><text x="490" y="76">Chrome</text>
<rect class="box alt" x="610" y="36" width="150" height="56" rx="10"/><text x="685" y="56">Agreement</text><text class="label" x="685" y="76">PDF</text>
<line class="edge" x1="170" y1="64" x2="206" y2="64"/>
<line class="edge" x1="370" y1="64" x2="406" y2="64"/>
<line class="edge" x1="570" y1="64" x2="606" y2="64"/>
</svg>

## Result

- RBI-compliant loan agreements in **10 languages**.
- A scalable, low-latency generation pipeline that runs inside the loan journey.
