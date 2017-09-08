---
title: Padding
status: Stable
---

Padding utilities are based on a global [spacing scale](/styleguide/css/styles/core/support/spacing) which helps keep horizontal and vertical spacing consistent. These utilities help us reduce the amount of custom CSS that could share the same properties, and allows to achieve many different page layouts using the same styles.

{:toc}

## Shorthand

Since padding utilities have many variations and will be used in many places, we use a shorthand naming convention to help keep class names succinct.

| Shorthand | Description |
| --- | --- |
| p | padding |
| t | top |
| r | right |
| b | bottom |
| l | left |
| x | horizontal, left & right |
| y | vertical, top & bottom |
| 0 | 0 |
| 1 | 4px |
| 2 | 8px |
| 3 | 16px |
| 4 | 24px |
| 5 | 32px |
| 6 | 40px |

_**Note:** The blue in the examples below represents the element, and the green represents the padding._

## Uniform spacing

Use uniform spacing utilities to apply equal padding to all sides of an element. These utilities can be used with a spacing scale from 0-6.

```html
<div class="padding-green d-inline-block p-0">
  <div class="d-inline-block block-blue">.p-0</div>
</div>
<div class="padding-green d-inline-block p-1">
  <div class="d-inline-block block-blue">.p-1</div>
</div>
<div class="padding-green d-inline-block p-2">
  <div class="d-inline-block block-blue">.p-2</div>
</div>
<div class="padding-green d-inline-block p-3">
  <div class="d-inline-block block-blue">.p-3</div>
</div>
<div class="padding-green d-inline-block p-4">
  <div class="d-inline-block block-blue">.p-4</div>
</div>
<div class="padding-green d-inline-block p-5">
  <div class="d-inline-block block-blue">.p-5</div>
</div>
<div class="padding-green d-inline-block p-6">
  <div class="d-inline-block block-blue">.p-6</div>
</div>
```

## Directional spacing

Use directional utilities to apply padding to an individual side, or the X and Y axis of an element. Directional utilities can change or override default padding, and can be used with a spacing scale of 0-6.

```html
<div class="padding-green d-inline-block pt-3">
  <div class="d-inline-block block-blue">.pt-3</div>
</div>
<div class="padding-green d-inline-block pr-3">
  <div class="d-inline-block block-blue">.pr-3</div>
</div>
<div class="padding-green d-inline-block pb-3">
  <div class="d-inline-block block-blue">.pb-3</div>
</div>
<div class="padding-green d-inline-block pl-3">
  <div class="d-inline-block block-blue">.pl-3</div>
</div>
<div class="padding-green d-inline-block py-3">
  <div class="d-inline-block block-blue">.py-3</div>
</div>
<div class="padding-green d-inline-block px-3">
  <div class="d-inline-block block-blue">.px-3</div>
</div>
```

## Responsive padding

All padding utilities can be adjusted per [breakpoint](/styleguide/css/styles/core/support/breakpoints) using the following formula: <br /> `p-[direction]-[breakpoint]-[spacer]`. Each responsive style is applied to the specified breakpoint and up.

```html
<div class="px-sm-2 px-lg-4 d-inline-block padding-green">
  <div class="d-inline-block block-blue">
    .px-sm-2 .px-lg-4
  </div>
</div>
```

## Responsive container padding

`.p-responsive` is a padding class that adds padding on the left and right sides of an element. On small screens, it gives the element padding of `$spacer-3`, on mid-sized screens it gives the element padding of `$spacer-6`, and on large screens, it gives the element padding of `$spacer-3`.

It is intended to be used with [container styles](/styleguide/css/styles/core/objects/grid#containers)

```html
<div class="container-lg p-responsive">
  <div class="border">
    .container-lg .p-responsive
  </div>
</div>
```
