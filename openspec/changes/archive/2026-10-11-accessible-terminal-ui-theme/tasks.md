## 1. Shared renderer

- [x] 1.1 Replace fixed dark panel surfaces with compact theme-inheriting field layout, semantic markers, multiline wrapping, and narrow-terminal behavior
- [x] 1.2 Replace filled table rows and headers with theme-neutral styles while retaining responsive complete-value rendering
- [x] 1.3 Keep plain output, TTY detection, sanitization, semantic text color, and QR contrast behavior unchanged

## 2. Verification

- [x] 2.1 Add exact renderer tests for no background/RGB escapes, all tones, multiline and long values, narrow layouts, tables, `NO_COLOR`, and sanitization
- [x] 2.2 Run focused terminal UI, report, application runtime, and strict OpenSpec tests
- [x] 2.3 Run the complete `make verify` release-candidate gate

## 3. Release

- [x] 3.1 Ship the archived change as the next patch and verify release surfaces, package channels, Pages, deleted feature branch, and clean synchronized main
