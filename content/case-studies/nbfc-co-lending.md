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
