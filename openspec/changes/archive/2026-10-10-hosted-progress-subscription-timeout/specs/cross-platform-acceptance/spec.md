## MODIFIED Requirements

### Requirement: Compiled hosted-route browser acceptance

Courier SHALL build the real executable before browser acceptance and SHALL serially exercise embedded hosted applications and HTTP routes in isolated homes with unique ports, mandatory process/resource cleanup, and a foreground idle interval beyond the ordinary IPC request timeout.

#### Scenario: Idle browser upload remains available

- **WHEN** a compiled foreground `from web://` delivery remains idle for more than five seconds before browser navigation
- **THEN** the initiating process remains active, the embedded upload application is reachable, and external stop exits cleanly with code `0`

#### Scenario: Browser uploads sequential files

- **WHEN** Playwright uploads Unicode text and binary files through a compiled `from web://` delivery
- **THEN** destination bytes match exactly, each committed name is present, and external stop makes the foreground command exit `0`

#### Scenario: Browser uploads an extraction archive

- **WHEN** Playwright uploads a safe tar.gz through `from web:// --extract`
- **THEN** the exact safe tree is committed without retaining the archive

#### Scenario: Browser downloads a directory

- **WHEN** Playwright opens `from <directory> to web://` and downloads the shared directory
- **THEN** the returned archive contains the expected nested byte-identical tree

#### Scenario: Administration runtime is exercised

- **WHEN** Playwright opens the embedded UI from a real foreground `ui start` and a second process runs `ui stop`
- **THEN** the page was reachable, the initiating process exits `0`, and administration state is removed
