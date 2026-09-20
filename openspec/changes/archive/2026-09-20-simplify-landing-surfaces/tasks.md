## 1. Landing presentation

- [x] 1.1 Remove the route metrics and section boundary rules, then verify the landing keeps exactly three semantic sections without the counter strip
- [x] 1.2 Replace installation tab styling with stable active and hover states, remove the repeated channel identity, and verify exact command copying remains unchanged
- [x] 1.3 Remove all CLI registry header, row, column, empty-state, and readout rules while preserving editable controls, focus visibility, and command construction

## 2. Acceptance and release readiness

- [x] 2.1 Update exact unit and Chromium acceptance assertions for the simplified surfaces and verify landing coverage remains 100%
- [x] 2.2 Inspect light and dark landing layouts at 320, 390, 1024, and 1440 pixels and verify stable tabs, no clipping, and no unexpected boundary rules
- [x] 2.3 Run strict OpenSpec validation, landing bundle checks, and `make verify` before archiving and shipping the change
