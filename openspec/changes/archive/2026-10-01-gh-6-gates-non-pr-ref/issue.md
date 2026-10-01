# Issue #6: `check-gates.sh` fails on a non-PR ref (e.g. `main`)

https://github.com/Kerman-Sanjuan/octospec/issues/6

## Summary

Running `scripts/check-gates.sh` on `main` fails G3 (`main` does not match `feat|fix/<issue>-<slug>`) and exits 1. The script assumes it is running on a PR branch.

## Steps to reproduce

1. Check out `main`.
2. `sh scripts/check-gates.sh`.

## Expected

G3 (and other PR-only gates) are skipped when `HEAD_REF` is a long-lived branch, or the script clearly reports "not a PR ref".

## Actual

`FAIL G3 branch name: 'main' does not match feat|fix/<issue>-<slug>` -> `FAIL gates failed`, exit 1.

## Impact

The script cannot be run on `main`; output is misleading.

## Fix

Skip G3 when `HEAD_REF` is `main`/`master`, or add an explicit PR-only mode.
