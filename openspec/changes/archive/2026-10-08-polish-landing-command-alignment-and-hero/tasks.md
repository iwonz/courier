## 1. Command surfaces

- [x] 1.1 Remove the visible shell-syntax label while retaining the localized accessible selector name, verified by landing component tests and a browser accessibility snapshot
- [x] 1.2 Normalize nested and direct readout padding so selector, command, and Copy share one inline-start axis, verified at desktop and mobile widths
- [x] 1.3 Restyle Copy as an icon-led text action with no hover background, verified by shared component tests and real-browser hover inspection

## 2. Hero composition

- [x] 2.1 Scale and position route-v3 as a responsive decorative layer behind the headline, verified at 390×844 and 1440×900
- [x] 2.2 Add theme-aware opacity and radial edge softening without modifying the raster asset, verified in light and dark themes
- [x] 2.3 Update project landing documentation and generated embedded CSS, verified by the web asset and embed checks

## 3. Verification

- [x] 3.1 Extend shared and landing regression coverage and retain exact 100% TypeScript/TSX coverage
- [x] 3.2 Run the production landing build and verify JavaScript and CSS budgets
- [x] 3.3 Run the complete browser acceptance suite and strict OpenSpec validation
