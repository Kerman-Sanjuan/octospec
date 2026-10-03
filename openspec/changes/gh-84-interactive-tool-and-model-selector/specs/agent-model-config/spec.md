## MODIFIED Requirements

### Requirement: An interactive TUI for the models
octospec SHALL provide an interactive TUI that sets the model of each role. The TUI SHALL offer a catalog of known models per tool, and SHALL allow a model outside the catalog to be entered by hand.

#### Scenario: Set a model in the TUI
- **WHEN** the user opens the model TUI
- **THEN** they pick a model for each role and the choice is saved to the configuration.

#### Scenario: Choose from the catalog
- **WHEN** the user opens a role
- **THEN** it lists the known models for the tools and the user selects one.

#### Scenario: Enter a model by hand
- **WHEN** the model the user wants is not in the catalog
- **THEN** the user can enter it by hand and it is saved to the configuration.
