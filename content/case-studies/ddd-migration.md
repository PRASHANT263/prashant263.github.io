---
title: From Monolith to Domain-Driven Design
summary: Restructuring tightly coupled lending systems into domain, service and gateway layers with clear boundaries.
order: 1
tags: [go, ddd, architecture]
metrics: ["Clear service boundaries", "Faster feature delivery", "Better maintainability"]
---
## Problem

The lending platform had grown as a tightly coupled monolith. Business rules for multiple lending products and partner integrations lived alongside database access and third-party calls, so a change in one flow could ripple into others. New features meant touching code far outside the area being changed.

## Approach

I led the migration to a Domain-Driven Design model built around three layers:

- **Domain layer**: the core lending concepts and rules, in plain Go, with no knowledge of databases, HTTP or partners.
- **Service layer**: use cases that orchestrate domain logic for a specific product flow.
- **Gateway layer**: adapters for everything external, such as databases, partner systems, queues and notifications.

The migration was incremental rather than a big-bang rewrite. Flows moved one at a time behind stable interfaces, so production traffic kept flowing throughout.

## Architecture

<svg class="diagram" viewBox="0 0 770 300" role="img" aria-label="Handlers call the service layer, which uses the domain layer and the gateway layer; gateways talk to the database, partner systems and queues">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="285" y="16" width="200" height="48" rx="10"/><text x="385" y="40">API handlers</text>
<rect class="box" x="285" y="110" width="200" height="56" rx="10"/><text x="385" y="130">Service layer</text><text class="label" x="385" y="150">product use cases</text>
<rect class="box" x="20" y="110" width="200" height="56" rx="10"/><text x="120" y="130">Domain layer</text><text class="label" x="120" y="150">lending rules, pure Go</text>
<rect class="box alt" x="550" y="110" width="200" height="56" rx="10"/><text x="650" y="130">Gateway layer</text><text class="label" x="650" y="150">adapters</text>
<rect class="box alt" x="420" y="232" width="105" height="44" rx="10"/><text x="472" y="254">MySQL</text>
<rect class="box alt" x="535" y="232" width="110" height="44" rx="10"/><text x="590" y="254">Partners</text>
<rect class="box alt" x="655" y="232" width="105" height="44" rx="10"/><text x="707" y="254">Queues</text>
<line class="edge" x1="385" y1="64" x2="385" y2="106"/>
<line class="edge" x1="285" y1="138" x2="224" y2="138"/>
<line class="edge" x1="485" y1="138" x2="546" y2="138"/>
<line class="edge" x1="630" y1="166" x2="480" y2="228"/>
<line class="edge" x1="645" y1="166" x2="595" y2="228"/>
<line class="edge" x1="665" y1="166" x2="700" y2="228"/>
</svg>

Dependencies point inward: the domain layer depends on nothing, and gateways implement interfaces defined by the services that use them.

## Result

- Clear service boundaries that make it obvious where new logic belongs.
- Better maintainability and faster feature delivery across lending products.
- Shared patterns reinforced through code reviews and design discussions with the engineers I mentor.
