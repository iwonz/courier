## MODIFIED Requirements

### Requirement: Observable browser upload lifecycle

The embedded browser-upload application SHALL accept one or more selected files, submit each file as an independent single-file transaction, and expose a current-batch list with each file's name, exact size, queued/uploading/delivered/failed state, and individual transfer progress. It SHALL derive active progress from browser upload events, claim delivery only after the server acknowledges the committed upload, continue remaining queued files after an individual rejection, prevent overlapping batch selections, and reset the file input after every completed batch so the same files can be selected again. Server failure text SHALL remain guarded and bounded, and batch summaries plus active progress SHALL use appropriate accessible live and progress semantics in every supported locale.

#### Scenario: Multiple selected files are delivered

- **WHEN** a visitor selects multiple files from an authorized browser-upload route
- **THEN** the workspace immediately lists every file and submits one request per file in deterministic selection order until each has an independent terminal result

#### Scenario: Active file reports progress

- **WHEN** the browser reports upload progress for the active file
- **THEN** only that file's progress indicator and percentage advance while later files remain visibly queued

#### Scenario: Selected file is committed

- **WHEN** an individual upload endpoint returns committed success
- **THEN** that file reaches 100 percent and is identified as delivered without removing other batch rows

#### Scenario: Selected file conflicts

- **WHEN** one selected file conflicts with an existing destination name
- **THEN** its row displays the guarded `409` collision failure and the next queued file is still attempted

#### Scenario: Completed batch is selected again

- **WHEN** every file in the current batch has succeeded or failed
- **THEN** the chooser is enabled, its native value is reset, and the same file selection can start a new batch

#### Scenario: Upload response is not a guarded JSON error

- **WHEN** an upload fails with a malformed, non-JSON, oversized, unsafe, or otherwise unusable error body
- **THEN** the affected row displays a bounded status-based failure without rendering response markup or protected data

