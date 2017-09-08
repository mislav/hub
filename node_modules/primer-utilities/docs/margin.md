---
title: Margin
status: Stable
---

Margin utilities are based on a global [spacing scale](/styleguide/css/styles/core/support/spacing) which helps keep horizontal and vertical spacing consistent. These utilities help us reduce the amount of custom CSS that share the same properties, and allows to achieve many different page layouts using the same styles.

{:toc}

## Naming convention

Since margin utilities have many variations and will be used in many places, we use a shorthand naming convention to help keep class names succinct.


| Shorthand | Description |
| --- | --- |
| m | margin |
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

_**Note:** The blue in the examples represents the element, and the orange represents the margin_

## Uniform spacing

Use uniform spacing utilities to apply equal margin to all sides of an element. These utilities can change or override default margins, and can be used with a spacing scale from 0-6.

```html
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-0">.m-0</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-1">.m-1</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-2">.m-2</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-3">.m-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-4">.m-4</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-5">.m-5</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue m-6">.m-6</div>
</div>
```

## Directional spacing

Use directional utilities to apply margin to an individual side, or the X and Y axis of an element. Directional utilities can change or override default margins, and can be used with a spacing scale of 0-6.

```html
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue mt-3">.mt-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue mr-3">.mr-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue mb-3">.mb-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue ml-3">.ml-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue my-3">.my-3</div>
</div>
<div class="margin-orange d-inline-block">
  <div class="d-inline-block block-blue mx-3">.mx-3</div>
</div>
```

## Center elements

Use `mx-auto`to center block elements with a set width.

```html
<div class="margin-orange">
  <div class="block-blue mx-auto text-center" style="width: 5rem;">.mx-auto</div>
</div>
```

## Reset margins
Reset margins built into typography elements or other components with `m-0`, `mt-0`, `mr-0`, `mb-0`, `ml-0`, `mx-0`, and `my-0`.

```html
<p class="mb-0 block-blue">No bottom margin on this paragraph.</p>
```

## Responsive margins

All margin utilities, except `mx-auto`, can be adjusted per [breakpoint](/styleguide/css/modules/grid#breakpoints) using the following formula: `m[direction]-[breakpoint]-[spacer]`. Each responsive style is applied to the specified breakpoint and up.

```html
<div class="d-inline-block margin-orange">
  <div class="mx-sm-2 mx-lg-4 d-inline-block block-blue">
    .mx-sm-2 .mx-lg-4
  </div>
</div>
```

## Negative margins

You can add negative margins to the top, right, bottom, or left of an item by adding a negative margin utility. The formula for this is: `m[direction]-n[spacer]`. This also works responsively, with the following formula: `m[direction]-[breakpoint]-n[spacer]`.

```html
<div class="d-inline-block margin-orange p-3">
  <div class="d-inline-block block-blue mt-n4 mr-lg-n4">
    .mt-n4 .mr-lg-n6
  </div>
</div>
```
