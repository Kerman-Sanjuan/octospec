# Issue #9: `/idea` and `/bug` cannot create issues from their YAML forms

https://github.com/Kerman-Sanjuan/octospec/issues/9

## Summary

`/idea` and `/bug` both instruct:

```
gh issue create --template feature.yml --body "<fields>"
gh issue create --template bug.yml --body "<fields>"
```

But `gh` rejects this: `--template is not supported when using --body or --body-file`. YAML issue forms cannot be pre-filled from `--body`.

## Steps to reproduce

1. `gh issue create --template feature.yml --title x --body y`
2. Observe the error.

## Expected

The documented primary command works, or the shim uses the body fallback by default.

## Actual

The primary path always errors; the commands only work via their written "if the form is unavailable" fallback.

## Impact

Confusing and fragile; every `/idea`/`/bug` run depends on the fallback.

## Fix

Make "pass the sections as the body" the documented primary path. Keep `feature.yml`/`bug.yml` as the human web-form contract (their `labels:` and required fields still matter for the UI).
