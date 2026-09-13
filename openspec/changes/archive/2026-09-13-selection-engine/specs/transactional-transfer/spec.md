## MODIFIED Requirements

### Requirement: Recursive non-destructive copy

Courier SHALL copy selected regular files, directories, and symlinks recursively into an absent final path while leaving source objects, excluded objects, and unrelated destination entries unchanged.

#### Scenario: Nested directory tree

- **WHEN** the source contains selected and excluded nested objects and the final path is absent
- **THEN** destination contains the selected supported object types, content, and structure, source still exists, and excluded and sibling destination entries remain unchanged
