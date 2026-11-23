# Research target

Verified against the official Postman pages. `scripts/research.sh` reads the `url:` line.

```yaml
url: https://learning.postman.com/flows/overview
product: Postman Flows
scope: palette, canvas, nodes, edges, run with live per-node status, persistence
out_of_scope: multi-user, cloud sync, accounts, sharing, imports, plugin marketplace, mobile, themes, i18n
auth: public docs need none; the live builder needs a Postman account — do not bypass
tos_notes: public docs and public product only; never copy target source, assets, CSS, icons, microcopy
```

## Official sources

| Source | URL | Use |
|---|---|---|
| Docs overview | https://learning.postman.com/flows/overview | what a flow is, start here |
| Build flows | https://learning.postman.com/flows/build-flows/overview | canvas, nodes, connections, run |
| Blocks reference | https://learning.postman.com/flows/reference/blocks/overview | node catalogue and config |
| Tutorials | https://learning.postman.com/flows/tutorials/overview | intended journeys |
| Product page | https://www.postman.com/product/flows/ | positioning and vocabulary |

Add the changelog, community feature requests, and public API/collection docs to
`spec/INTEL.md` as they are found. Do not guess a URL; if a page is unreachable, say so.

## Jobs to be done

1. Open the builder and see the node palette.
2. Drag a node onto the canvas.
3. Connect a declared output port to an input port.
4. Configure a node and reopen it.
5. Run the flow and watch each node's status.
6. Make a node fail and see descendants skip.
7. Move a node, restart, and find it where it was.

## Scope guard

The clone is the surface in `GOAL.md`: palette, canvas, nodes, edges, run/status, SQLite
persistence. Postman-specific concepts such as Blocks, FQL, environments, and collections are
intel for `spec/INTEL.md`; they are not automatically in scope. Anything outside `GOAL.md`
goes to `spec/CONTRACT-PROPOSAL.md` or `docs/SCALE.md`, not into the clone.

## Auth and data guard

- Public docs and the public product page need no login.
- The live builder needs a Postman account. Use a dedicated account and only your own data.
- Never bypass authentication, scrape another user's data, or copy target source, assets, CSS,
  icons, class names, or microcopy.
- The target's shipped JavaScript and sourcemaps stay off-limits. See `docs/RESEARCH.md`.
