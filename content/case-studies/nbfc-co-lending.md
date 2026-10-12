---
title: Co-Lending Integrations with 4 NBFC Partners
summary: End-to-end integrations with Tata Capital, OFL, ABCL and IFMR that grew co-lending loan volume by 20%.
order: 2
tags: [go, integrations, fintech]
metrics: ["+20% loan volume", "4 NBFC partners", "RBI-compliant"]
---
## Context

Co-lending lets a lending platform originate loans together with NBFC partners. Each partner added widens the platform's lending distribution network, and every integration has to meet RBI requirements.

## What I did

- Integrated **Tata Capital, OFL, ABCL and IFMR** end to end, expanding the platform's lending distribution network.
- Designed and shipped the **Multi-Lender Display** feature, owning the server-side business logic that lets users compare and select their preferred lender within a single loan journey.
- Handled production debugging of critical fund-flow issues, tracing root causes across backend services.
- Collaborated with product, finance, compliance and operations teams to ship RBI-compliant features and circulars on schedule.

## Architecture

<svg class="diagram" viewBox="0 0 770 230" role="img" aria-label="The lending platform's loan journey connects to four NBFC co-lending partners">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="235" y="16" width="300" height="56" rx="10"/><text x="385" y="36">Lending platform</text><text class="label" x="385" y="56">loan journey · lender selection</text>
<rect class="box alt" x="20" y="160" width="160" height="48" rx="10"/><text x="100" y="184">Tata Capital</text>
<rect class="box alt" x="210" y="160" width="160" height="48" rx="10"/><text x="290" y="184">OFL</text>
<rect class="box alt" x="400" y="160" width="160" height="48" rx="10"/><text x="480" y="184">ABCL</text>
<rect class="box alt" x="590" y="160" width="160" height="48" rx="10"/><text x="670" y="184">IFMR</text>
<line class="edge" x1="330" y1="72" x2="110" y2="156"/>
<line class="edge" x1="370" y1="72" x2="292" y2="156"/>
<line class="edge" x1="400" y1="72" x2="478" y2="156"/>
<line class="edge" x1="440" y1="72" x2="660" y2="156"/>
</svg>

A simplified view: one loan journey, four NBFC partners integrated end to end.

## Result

- **+20% co-lending loan volume** from the expanded partner network.
- Users can compare lenders and pick their preferred one inside a single loan journey.
