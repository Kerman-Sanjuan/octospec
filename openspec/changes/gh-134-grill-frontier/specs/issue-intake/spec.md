## ADDED Requirements

### Requirement: The intake grills the decision frontier before filing

`/idea` and `/bug` SHALL treat the request as a decision tree, ask in rounds only the decisions whose
prerequisites are settled, and SHALL NOT file the issue while any decision on the frontier is unanswered.

#### Scenario: A decision depends on one still open

- **WHEN** a question's answer depends on a question asked in the same round
- **THEN** the dependent question is deferred to the next round instead of being asked twice

#### Scenario: The frontier is empty

- **WHEN** every branch of the decision tree has an answer
- **THEN** the agent files the issue and stops

#### Scenario: The frontier is not empty

- **WHEN** the user stops answering with a decision still open
- **THEN** no issue is filed and the open decision is named to the user

### Requirement: Every question carries a recommendation

Each question in a round SHALL be numbered and SHALL state the agent's recommended answer, worded so
that a plain "yes" accepts it.

#### Scenario: A round of questions

- **WHEN** the agent asks three questions in one round
- **THEN** each is numbered and each has a recommended answer beneath it

### Requirement: Facts are looked up, decisions are asked

The agent SHALL look up every fact it can read from the repository, the filesystem, and the issue
tracker, and SHALL ask the user only for decisions.

#### Scenario: A question needs a fact from the repository

- **WHEN** a frontier question can be answered by reading the code or an existing issue
- **THEN** the agent reads it and states the answer instead of asking

### Requirement: The frontier covers the seven groups

The frontier SHALL cover who acts, the observable outcome, what done looks like, the failures that
matter, the boundary and why it is out of scope, the constraints already known, and the evidence that
will prove the outcome.

#### Scenario: A group has no answer

- **WHEN** the boundary of the change is undecided
- **THEN** the frontier still holds that question and the issue is not filed

### Requirement: The answers are recorded in the issue

`/idea` and `/bug` SHALL write the closed decisions into the issue body under a `## Grill` section, one
decision per line, and `/spec` SHALL read that section and SHALL NOT re-ask a question it answers.

#### Scenario: A filed issue

- **WHEN** the grill ends
- **THEN** the issue body carries `## Grill` with one line per closed decision and no open question

#### Scenario: A later spec

- **WHEN** `/spec` turns an issue into a change
- **THEN** it reads `## Grill` and asks nothing already answered there

### Requirement: An unattended grill stops instead of inventing

When the grill runs with no human available, it SHALL NOT invent an answer: it SHALL stop with
`reason: needs_human`, naming the question that is blocking.

#### Scenario: Unattended run with an open decision

- **WHEN** the grill runs unattended and a frontier question has no recorded answer
- **THEN** the stage stops with `reason: needs_human` and names the question
