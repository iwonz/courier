# Design: Borderless UI flow

## Continuous canvas

Courier browser applications use the document background as their only page-scale surface. Sections and workflows are arranged directly on that canvas. Large rounded rectangles, translucent boards, card shadows, and nested tinted groups do not define ordinary hierarchy.

Spacing, type scale, alignment, and a small number of hairline dividers communicate structure. Dividers are permitted between independently scrollable columns, manifest rows, and navigation/inspection regions where the boundary explains behavior. Form fields, alerts, selected controls, focus rings, status badges, and destructive actions retain explicit chrome because it communicates state.

## Shared primitives

The shared Card primitive remains available as a semantic layout wrapper but has no default background, radius, shadow, border, or blur. CommandReadout and RouteDisplay become transparent inline workflow structures. Tabs use a line-based selected state instead of a containing pill surface.

## Surface composition

The landing route composition, install command, and CLI registry sit directly in normal flow. Delivery authentication and manifests no longer sit inside cards. Administration counters, navigator, inspector, route facts, and policy controls share the document canvas without rounded panels. Responsive behavior, scroll bounds, protected-data isolation, and keyboard semantics stay unchanged.

Landing sections use compact adjacent spacing instead of stacking large top and bottom padding. The masthead is fixed to the visual viewport with no background fill or blur; the document reserves its height so content is never covered. Selected route controls retain the filled primary treatment and are never overridden by a translucent inactive-control background.
