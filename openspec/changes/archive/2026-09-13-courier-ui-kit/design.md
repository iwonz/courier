# Context

The supplied HTML files are visual references rather than source code. They contain embedded image data, obsolete mirror-oriented copy, and browser-injected AdGuard scripts. Future Courier browser surfaces need a small shared library that can be built deterministically and embedded without retaining unrelated page code.

# Decisions

## One framework-neutral package boundary with Lit components

The `@courier/ui` package owns tokens, icons, localization, theme state, and reusable Lit components. Applications import its public entry point and catalogs rather than copying CSS or component implementations. The package is built as an ES module with Vite and emits declarations through TypeScript.

## Identity assets are explicit and auditable

Only selected Courier-owned raster or vector assets are copied into a dedicated asset directory. A machine-readable provenance manifest records the source reference, media type, digest, and intended use. No executable script, remote resource, page markup, obsolete mirror message, or credentials are retained.

## Theme and locale are typed state

Theme preference is the union `system | light | dark`, persists in local storage when available, and resolves system changes through `matchMedia`. Localization keys are derived from the English catalog; every locale must implement the same shape. English is the fallback, browser languages choose the initial locale, and adding a catalog does not require backend changes.

## Quality gates operate on source

Vitest with V8 coverage enforces 100% statements, branches, functions, and lines for first-party TypeScript. Vite library builds and an asset-integrity checker run from `make verify`. Generated output and dependencies are excluded from coverage and source control.

# Risks / Trade-offs

The task intentionally creates primitives rather than complete delivery or administration pages. Browser end-to-end matrices are added with those consuming applications, while this change tests DOM behavior in a deterministic local environment.

# Migration Plan

Create the package and provenance manifest, implement and test tokens/theme/i18n/components, wire build and coverage into the repository gate, verify extracted assets, delete the reference directory, then archive this OpenSpec change.
