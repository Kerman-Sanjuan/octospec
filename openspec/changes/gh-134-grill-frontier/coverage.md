# Coverage — gh-134-grill-frontier

Each requirement in this change, and the thing that proves it.

- The intake grills the decision frontier before filing -> TestIdeaCommandGrillsTheFrontier
- Every question carries a recommendation -> TestIdeaCommandRecommendsAnAnswer
- Facts are looked up, decisions are asked -> TestIdeaAgentLooksFactsUpAndAsksDecisions
- The frontier covers the seven groups -> TestIdeaAgentCoversTheSevenGroups
- The answers are recorded in the issue -> TestSpecCommandReadsTheGrill
- An unattended grill stops instead of inventing -> TestIdeaAgentStopsUnattended
- G1 issue sections -> scripts/check-gates.sh
