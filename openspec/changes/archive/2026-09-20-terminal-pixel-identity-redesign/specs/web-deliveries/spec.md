## MODIFIED Requirements

### Requirement: Pixel delivery identity

The embedded delivery React application SHALL render as a terminal transfer console using mark-v3 and delivery-v3. It SHALL separate title, real route/status, authorization, and manifest regions with flat surfaces and one-pixel rules. Manifest entries SHALL be ruled rows. Authorization, upload, download, navigation, protected metadata isolation, and API-only behavior SHALL remain unchanged, and no fictional transfer data SHALL be rendered.

#### Scenario: An authorized manifest renders

- **WHEN** the runtime API returns delivery metadata
- **THEN** the console shows the actual route state and ruled manifest while preserving the same upload, download, and directory navigation actions

#### Scenario: A visitor opens a delivery

- **WHEN** the visitor opens a protected or authorized delivery route
- **THEN** mark-v3 and delivery-v3 orient the transfer console while authentication and protected metadata behavior remain authoritative
