#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
SCRIPT="$ROOT/scripts/prepare-changelog.sh"
TMP="$(mktemp -d)"
trap 'rm -rf "$TMP"' EXIT

fail() {
	echo "FAIL: $*" >&2
	exit 1
}

# Missing Unreleased
printf '# Changelog\n\n## [0.0.1] - 2026-01-01\n' >"$TMP/empty.md"
if CHANGELOG_FILE="$TMP/empty.md" "$SCRIPT" 0.1.0 --check-only; then
	fail "expected check-only to fail without [Unreleased]"
fi

# Unreleased without ### notes
printf '# Changelog\n\n## [Unreleased]\n\n## [0.0.1] - 2026-01-01\n' >"$TMP/blank.md"
if CHANGELOG_FILE="$TMP/blank.md" "$SCRIPT" 0.1.0 --check-only; then
	fail "expected check-only to fail without ### notes"
fi

cat >"$TMP/CHANGELOG.md" <<'EOF'
# Changelog

## [Unreleased]

### Добавлено

- first public API

## [0.0.1] - 2026-01-01
EOF

CHANGELOG_FILE="$TMP/CHANGELOG.md" "$SCRIPT" 0.1.0 --check-only >/dev/null
CHANGELOG_FILE="$TMP/CHANGELOG.md" "$SCRIPT" 0.1.0 2026-08-13 >/dev/null

grep -q '^## \[Unreleased\]$' "$TMP/CHANGELOG.md" || fail "Unreleased header missing after promote"
grep -q '^## \[0.1.0\] - 2026-08-13$' "$TMP/CHANGELOG.md" || fail "version header missing"
grep -q 'first public API' "$TMP/CHANGELOG.md" || fail "notes were dropped"

# Idempotent when the version section already exists
CHANGELOG_FILE="$TMP/CHANGELOG.md" "$SCRIPT" 0.1.0 --check-only >/dev/null
CHANGELOG_FILE="$TMP/CHANGELOG.md" "$SCRIPT" 0.1.0 2026-08-14 >/dev/null
count="$(grep -c '^## \[0.1.0\]' "$TMP/CHANGELOG.md")"
[ "$count" = "1" ] || fail "expected a single 0.1.0 section, got $count"

echo "prepare-changelog.sh self-test ok"
