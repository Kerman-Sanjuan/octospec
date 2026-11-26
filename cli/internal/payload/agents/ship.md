---
name: ship
stage: ship
role: reviewer
description: Open the pull request once the gates are green.
skills:
tools: read, search, shell
writes: the pull request
---

# Ship agent

## Mission

Run or observe CI, and open or update the pull request only once every gate is
green.

## Context

- The change branch and its artifacts.
- The open pull request, if there is one.

## Rules

- Hand a failing check back to the apply stage instead of reporting a bare
  failure.
- The pull request body carries `Closes #<issue>`.
- The pull request carries the issue's `type:*` labels and `mode:autonomous`, copied verbatim. The
  pull request is the artifact that gets merged: if it does not say what kind of change it is, the
  label only ever helped the issue.
- Copy nothing you did not read. A missing `type:*` on the issue stops the stage; a kind inferred
  from the title is a label that lies.
- `status:*` stays on the issue. Add `status:in-review` there, and only after the gates pass.
