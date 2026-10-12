---
title: From Monolith to Domain-Driven Design
summary: Leading the migration of lending systems from a tightly coupled monolith to domain, service and gateway layers with clear boundaries.
order: 1
tags: [go, ddd, architecture]
metrics: ["Clear service boundaries", "Faster feature delivery", "Better maintainability"]
---
## Context

The lending systems had grown as a tightly coupled monolith, supporting multiple lending products and partner integrations. Tight coupling makes service boundaries unclear and slows down maintenance and feature delivery.

## What I did

As a Senior Software Engineer, I led the migration of these lending systems to a **Domain-Driven Design** model, introducing three dedicated layers:

- **Domain layer**: the core lending concepts and business rules.
- **Service layer**: the use cases that orchestrate those rules for each product flow.
- **Gateway layer**: the integration points with everything external, such as databases, partners and messaging.

Alongside the migration, I manage architectural improvements across critical lending workflows, and I reinforce the patterns through code reviews and design discussions with the five engineers I mentor.

## Architecture

<svg class="diagram" viewBox="0 0 770 300" role="img" aria-label="Handlers call the service layer, which uses the domain layer and the gateway layer; gateways connect to databases, partners and messaging">
<defs><marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse"><path d="M0 0L10 5L0 10z"/></marker></defs>
<rect class="box" x="285" y="16" width="200" height="48" rx="10"/><text x="385" y="40">API handlers</text>
<rect class="box" x="285" y="110" width="200" height="56" rx="10"/><text x="385" y="130">Service layer</text><text class="label" x="385" y="150">product use cases</text>
<rect class="box" x="20" y="110" width="200" height="56" rx="10"/><text x="120" y="130">Domain layer</text><text class="label" x="120" y="150">lending rules</text>
<rect class="box alt" x="550" y="110" width="200" height="56" rx="10"/><text x="650" y="130">Gateway layer</text><text class="label" x="650" y="150">integrations</text>
<rect class="box alt" x="420" y="232" width="105" height="44" rx="10"/><text x="472" y="254">Database</text>
<rect class="box alt" x="535" y="232" width="110" height="44" rx="10"/><text x="590" y="254">Partners</text>
<rect class="box alt" x="655" y="232" width="105" height="44" rx="10"/><text x="707" y="254">Messaging</text>
<line class="edge" x1="385" y1="64" x2="385" y2="106"/>
<line class="edge" x1="285" y1="138" x2="224" y2="138"/>
<line class="edge" x1="485" y1="138" x2="546" y2="138"/>
<line class="edge" x1="630" y1="166" x2="480" y2="228"/>
<line class="edge" x1="645" y1="166" x2="595" y2="228"/>
<line class="edge" x1="665" y1="166" x2="700" y2="228"/>
</svg>

A simplified view of the target layering.

## Result

- Clear service boundaries between domain, service and gateway code.
- Improved maintainability and feature delivery across lending products.
- Better production reliability across critical lending workflows.
