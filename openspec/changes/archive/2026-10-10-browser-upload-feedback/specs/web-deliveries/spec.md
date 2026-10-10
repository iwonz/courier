## ADDED Requirements

### Requirement: Observable browser upload lifecycle

The embedded browser-upload application SHALL expose pending, committed-success, and rejected-failure states without discarding authorized delivery metadata. It SHALL claim success only after the server acknowledges the committed upload, SHALL render guarded HTTP failure status and reason as text, SHALL prevent overlapping selections, and SHALL reset the file input after every completed attempt so the same file can be selected again. Pending and completed outcomes SHALL be announced with appropriate accessible live semantics in every supported locale.

#### Scenario: Selected file is committed

- **WHEN** a visitor selects a file and the upload endpoint returns committed success
- **THEN** the chooser shows an indeterminate pending state during the request and then identifies the successfully delivered file

#### Scenario: Selected file conflicts

- **WHEN** a visitor selects a file whose destination name already exists
- **THEN** the page retains the delivery workspace, displays the guarded `409` collision failure, and permits the same file to be selected again

#### Scenario: Upload response is not a guarded JSON error

- **WHEN** an upload fails with a malformed, non-JSON, or otherwise unusable error body
- **THEN** the page displays a bounded status-based failure without rendering response markup or protected data
