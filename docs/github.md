# GitHub repository page

This checklist is for maintainers. Files in `.github/` apply on merge; **About**, topics, and feature toggles still need a one-time pass in the GitHub UI (**Settings** and the gear next to About).

Current snapshot (before this pass): empty description, no topics, default labels only, community health ~37%, license detected as “Other”, wiki and Projects on, no Discussions, no GitHub Release yet.

## About (gear on the repo header)

Paste these values:

| Field | Value |
| --- | --- |
| Description | `Rule-based Russian word and sentence tokenizer. Go port of natasha/razdel with behavioral parity.` |
| Website | `https://pkg.go.dev/github.com/muonsoft/go-razdel` |
| Topics | see list below |
| Releases | leave checked (badge in README) |
| Packages | optional |
| Deployments | off unless used |

Topics (GitHub allows up to 20; this set matches other `muonsoft` libraries plus NLP discoverability):

```
go
golang
golang-package
nlp
tokenizer
russian
text-segmentation
sentence-segmentation
razdel
natural-language-processing
```

Do not put Linear/internal ticket IDs in the public description.

## Features (Settings → General)

| Feature | Recommendation | Why |
| --- | --- | --- |
| Issues | **On** | Already on; templates live in `.github/ISSUE_TEMPLATE/`. |
| Discussions | Off for now | Small library; issues are enough. Revisit if usage questions pile up. |
| Wiki | **Off** | Stale duplicate of README/`docs/`. |
| Projects | **Off** unless the org uses them | Empty Projects tab adds noise. |
| Sponsors / Funding | Off unless there is a funding file | — |
| Preserve this repository | org policy | — |

## Settings → General (other)

- Default branch: `main`.
- Allow squash merge (keep history readable). Rebase/merge: org preference.
- Automatically delete head branches: on.
- Template repository: no.

## Settings → Code security

- **Private vulnerability reporting**: on.
- Dependabot alerts: on (Go module is tiny; still useful for Actions pins).

A short policy is in [`SECURITY.md`](../SECURITY.md).

## Settings → Pages

Leave **GitHub Pages** off. API docs are on [pkg.go.dev](https://pkg.go.dev/github.com/muonsoft/go-razdel); narrative docs are in this repo.

## License detection

`LICENSE` is canonical MIT so GitHub can show the MIT badge (previously “Other” because upstream attribution sat in the same file). Third-party notice is [`NOTICE`](../NOTICE).

## Labels

Defined in [`.github/labels.yml`](../.github/labels.yml). After merge to `main`, workflow **Label sync** creates or updates them (`workflow_dispatch` can run it immediately). Existing labels are not deleted (`skip-delete: true`).

### Keep GitHub defaults

| Label | Use for |
| --- | --- |
| `bug` | Incorrect tokenize/sentenize result, panic, broken invariant |
| `enhancement` | New behavior or API that is not a parity fix |
| `documentation` | README, contracts, examples |
| `question` | Usage; close or convert if it becomes a bug |
| `good first issue` | Isolated, documented, no submodule bump |
| `help wanted` | Maintainer agrees it is extra-hands work |
| `duplicate` / `invalid` / `wontfix` | Triage |

### Project labels

| Label | Use for |
| --- | --- |
| `tokenize` | `Tokenize` / atom / join rules |
| `sentenize` | `Sentenize` / abbreviations / bounds |
| `parity` | Mismatch or intentional mismatch vs Python `razdel` |
| `performance` | Allocations, benchmarks, hot path |
| `ci` | Workflows, test harness, fixtures |
| `api` | Public types, offset contract, module versioning |
| `upstream` | `third_party/razdel` pin, natasha/razdel issues |

Combine freely: a smile-tokenization bug is `bug` + `tokenize` + `parity`.

Issue forms only auto-apply labels that already exist on GitHub (`bug`, `enhancement`, `documentation`). Add `tokenize` / `sentenize` / `parity` while triaging.

## Community files (in git)

| File | Role |
| --- | --- |
| [`README.md`](../README.md) | Landing page for users |
| [`CONTRIBUTING.md`](../CONTRIBUTING.md) | How to build and test |
| [`SECURITY.md`](../SECURITY.md) | How to report a vulnerability |
| [`.github/ISSUE_TEMPLATE/`](../.github/ISSUE_TEMPLATE/) | Bug, feature, parity forms |
| [`.github/PULL_REQUEST_TEMPLATE.md`](../.github/PULL_REQUEST_TEMPLATE.md) | PR checklist |
| [`.github/labels.yml`](../.github/labels.yml) | Label source of truth |

## Releases and social preview

1. Publish the first GitHub Release (`v0.1.0`) via **Actions → Release** so the README release badge is not empty.
2. Optional: **Settings → General → Social preview** — a simple image with the name `go-razdel` and the one-line description. Skip until a designer pass; the README is enough for v0.

## What this repo will not host

- A second copy of the Python README evaluation tables (they belong to [natasha/razdel](https://github.com/natasha/razdel) and [naeval](https://github.com/natasha/naeval)).
- GitHub Pages marketing site.
- Discussions until there is real question volume.
