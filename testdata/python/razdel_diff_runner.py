#!/usr/bin/env python3
"""Batch JSON driver for Go differential tests: compares token/sentence texts with upstream razdel.

stdin:  {"cases":[{"mode":"tokenize"|"sentenize","text":"..."}, ...]}
stdout: {"results":[["tok1",...], ...]} — one string list per case, same order.

Requires PYTHONPATH pointing at the upstream repo root (directory containing the `razdel` package).
"""
import json
import sys

from razdel import sentenize, tokenize


def _run_one(mode, text):
    if mode == "tokenize":
        return [t.text for t in tokenize(text)]
    if mode == "sentenize":
        return [s.text for s in sentenize(text)]
    raise ValueError("unknown mode %r" % (mode,))


def main():
    req = json.load(sys.stdin)
    cases = req["cases"]
    results = [_run_one(c["mode"], c["text"]) for c in cases]
    json.dump({"results": results}, sys.stdout, ensure_ascii=False)
    sys.stdout.write("\n")


if __name__ == "__main__":
    main()
