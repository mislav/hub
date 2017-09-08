---
title: Breakpoints
status: Stable
source: https://github.com/primer/primer-css/blob/master/modules/primer-support/lib/mixins/layout.scss
---

{:toc}

Our breakpoints are based on screen widths where our content starts to break. Since most of GitHub is currently a fixed-width with we use pixel widths to make it easy for us to match page widths for responsive and non-responsive pages. **Our breakpoints may change as we move more of the product into responsive layouts.**

We use abbreviations for each breakpoint to keep the class names concise. This abbreviated syntax is used consistently across responsive styles. Responsive styles allow you to change the styles properties at each breakpoint. For example, when using column widths for grid layouts, you can change specify that the width is 12 columns wide at the small breakpoint, and 6 columns wide from the large breakpoint: `<div class="col-sm-12 col-lg-6">...</div>`

| Breakpoint | Syntax | Description |
| --- | --- | --- |
| Small | sm | min-width: 544px |
| Medium | md | min-width: 768px |
| Large | lg | min-width: 1012px |
| Extra-large | xl | min-width: 1280px |

<small>**Note:** The `lg` breakpoint matches our current page width of `980px` including left and right padding of `16px` (`$spacer-3`). This is so that content doesn't touch the edges of the window when resized.</small>

Responsive styles are available for [margin](./utilities/margin#responsive-margin), [padding](./utilities/padding#responsive-padding), [layout](./utilities/layout), [flexbox](.utilities/flexbox#responsive-flex-utilities), and the [grid](./objects/grid#responsive-grids) system.

## Breakpoint variables

The above values are defined as variables, and then put into a Sass map that generates the media query mixin.

```scss

// breakpoints
$width-xs: 0;
$width-sm: 544px;
$width-md: 768px;
$width-lg: 1012px;
$width-xl: 1280px;

$breakpoints: (
  // Small screen / phone
  sm: $width-sm,
  // Medium screen / tablet
  md: $width-md,
  // Large screen / desktop (980 + (12 * 2)) <= container + gutters
  lg: $width-lg,
  // Extra large screen / wide desktop
  xl: $width-xl
) !default;

```

## Media query mixins
Use media query mixins when you want to change CSS properties at a particular breakpoint. The media query mixin works by passing in a breakpoint value, such as `breakpoint(md)`.

Media queries are scoped from each breakpoint and upwards. In the example below, the font size is `28px` until the viewport size meets the `lg` breakpoint, from there upwards—including through the `xl` breakpoint—the font size will be `32px`.

```
.styles {
  font-size: 28px;
  @include breakpoint(md) { font-size: 32px; }
}
```
