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
