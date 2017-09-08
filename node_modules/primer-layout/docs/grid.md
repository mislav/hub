---
title: Grid
status: New release
status_issue: https://github.com/github/design-systems/issues/88
source: https://github.com/primer/primer-css/blob/master/modules/primer-layout/lib/grid.scss
---

The grid is 12 columns and percentage-based. The number of columns a container spans can be adjusted across breakpoints for responsive layouts. The grid system works with a variety of layout utilities to achieve different results.

{:toc}

## Float based grid

Use `.clearfix` on the container and float utilities with columns for a floated grid layout.

```html
<div class="container-lg clearfix">
  <div class="col-4 float-left border p-4">
    My column
  </div>
  <div class="col-4 float-left border p-4">
    Looks better
  </div>
  <div class="col-4 float-left border p-4">
    Than your column
  </div>
</div>
```

### Reversed grid

To reverse the order of columns, use `float-right` to float columns to the right.

```html
<div class="container-lg clearfix">
  <div class="col-4 float-right border p-4">
    One
  </div>
  <div class="col-4 float-right border p-4">
    Two
  </div>
  <div class="col-4 float-right border p-4">
    Three
  </div>
</div>
```

## Nesting
You can infinitely nest grid layouts within other columns since the column widths are percentage based. With great flexibility comes great responsibility - be sensible with how far you nest!

```html
<div class="clearfix">
  <div class="col-6 float-left px-1">
    <div class="border p-1">Unnested</div>
  </div>
  <div class="col-6 float-left">
    <div class="clearfix">
      <div class="col-6 float-left px-1">
        <div class="border p-1">1 x Nested</div>
      </div>
      <div class="col-6 float-left">
        <div class="col-6 float-left px-1">
          <div class="border p-1">2 x Nested</div>
        </div>
        <div class="col-6 float-left px-1">
          <div class="border p-1">2 x Nested</div>
        </div>
      </div>
    </div>
  </div>
</div>
```

## Centering a column

Use `.mx-auto` to center columns within a container.
```html
<div class="border">
  <div class="col-6 p-2 mx-auto border">
    This column is the center of my world.
  </div>
</div>
```


## Column widths
Column widths can be used with any other block or inline-block elements to add percentage-based widths.
```html
<div>
  <div class="col-4 float-right p-2 border text-red">
    <%= octicon "heart" %> Don't go bacon my heart.
  </div>
  <p>T-bone drumstick alcatra ribeye. Strip steak chuck andouille tenderloin bacon tri-tip ball tip beef capicola rump. Meatloaf bresaola drumstick ball tip salami. Drumstick ham bacon alcatra pig porchetta, spare ribs leberkas pork belly.</p>
</div>
```

## Offset columns

Using column offset classes can push a div over X number of columns. They work responsively using the [breakpoints outlined below](/styleguide/css/modules/grid#responsive-grids).

```html
<div class="clearfix">
  <div class="offset-1 col-3 border p-3">.offset-1</div>
  <div class="offset-2 col-3 border p-3">.offset-2</div>
</div>
```

## Gutters
Use gutter styles or padding utilities to create gutters. You can use the default gutter style, `gutter`, or either of its modifiers, `gutter-condensed` or `gutter-spacious`. Gutter styles also support responsive breakpoint modifiers. Gutter styles add padding to the left and right side of each column and apply a negative margin to the container to ensure content inside each column lines up with content outside of the grid.

```html
<div class="clearfix gutter-md-spacious border">
  <div class="col-3 float-left">
    <div class="border p-3">.col-md-3</div>
  </div>
  <div class="col-3 float-left">
    <div class="border p-3">.col-md-3</div>
  </div>
  <div class="col-3 float-left">
    <div class="border p-3">.col-md-3</div>
  </div>
  <div class="col-3 float-left">
    <div class="border p-3">.col-md-3</div>
  </div>
</div>
```

Use padding utilities to create gutters for more customized layouts.

```html
<div class="container-lg clearfix">
  <div class="col-3 float-left pr-2 mb-3">
    <div class="border bg-gray-light">.pr-2</div>
  </div>
  <div class="col-3 float-left px-2 mb-3">
    <div class="border bg-gray-light">.px-2</div>
  </div>
  <div class="col-3 float-left px-2 mb-3">
    <div class="border bg-gray-light">.px-2</div>
  </div>
  <div class="col-3 float-left pl-2 mb-3">
    <div class="border bg-gray-light">.pl-2</div>
  </div>
</div>
<div class="container-lg clearfix">
  <div class="col-3 float-left pr-2">
    <div class="border bg-gray-light">.pr-2</div>
  </div>
  <div class="col-9 float-left pl-2">
    <div class="border bg-gray-light">.pl-2</div>
  </div>
</div>
```


## Inline-block grids
Use column widths with `d-inline-block` as an alternative to floated grids.

```html
<div>
  <div class="col-4 d-inline-block border">
    .col-4 .d-inline-block
  </div><!--
  --><div class="col-4 d-inline-block border">
    .col-4 .d-inline-block
  </div><!--
  --><div class="col-4 d-inline-block border">
    .col-4 .d-inline-block
  </div>
</div>
```

You can use column widths and other utilities on elements such as lists to create the layout you need while keeping the markup semantically correct.
```html
<ul class="list-style-none">
  <li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/broccolini.png" alt="broccolini" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/jonrohan.png" alt="jonrohan" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/muan.png" alt="muan" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/pmarsceill.png" alt="pmarsceill" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/sophshep.png" alt="sophshep" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/cmwinters.png" alt="cmwinters" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/jeejkang.png" alt="jeejkang" /></li><!--
  --><li class="d-inline-block col-2 p-2"><img class="width-full avatar" src="/mdo.png" alt="mdo" /></li>
</ul>
```


## Display table grids
Using [display table utilities](/styleguide/css/utilities/layout#display) with columns gives you some alternative layout options.

A useful example is being able to keep the height of the container equal across a row when the length of content may differ.

```html
<div class="d-table col-12">
  <div class="col-4 d-table-cell border p-2">
    Bacon ipsum dolor amet leberkas pork pig kielbasa shankle ribeye meatball, salami alcatra venison.
  </div><!--
  --><div class="col-4 d-table-cell border p-2">
    Pork chop cupim cow turkey frankfurter, landjaeger fatback hamburger meatball salami spare ribs. Rump tenderloin salami, hamburger frankfurter landjaeger andouille.
  </div><!--
  --><div class="col-4 d-table-cell border p-2">
    Brisket tongue frankfurter cupim strip steak rump picanha pancetta pork pig kevin pastrami biltong. Shankle venison meatball swine sausage ground round. Tail pork loin ribeye kielbasa short ribs pork chop.
  </div>
</div>
```
You can also create an alternative [media object](/styleguide/css/utilities/layout#the-media-object) layout with `.display-table` and column widths.

```html
<div class="d-table col-12">
  <div class="col-2 d-table-cell v-align-middle">
    <img class="width-full avatar" src="/github.png" alt="github" />
  </div>
  <div class="col-10 d-table-cell v-align-middle pl-4">
    <h1 class="text-normal lh-condensed">GitHub</h1>
    <p class="h4 text-gray text-normal mb-2">How people build software.</p>
    <a class="text-gray text-small" href="#url">https://github.com/about</a>
  </div>
</div>
```

Note that table cells will fill the width of their container even when the total columns doesn't add up to 12.

```html
<div class="d-table col-12">
  <div class="col-4 d-table-cell border">
    .col-4 .d-table-cell
  </div><!--
  --><div class="col-4 d-table-cell border">
    .col-4 .d-table-cell
  </div><!--
  --><div class="col-2 d-table-cell border">
    .col-2 .d-table-cell
  </div>
</div>
```

## Flexbox grids

You can use [flex utilities](/styleguide/css/utilities/flexbox) on the container and columns to create a flexbox grid.

This can be useful for keeping columns the same height, justifying content and vertically aligning items. The flexbox grid is also great for working with responsive layouts.

```html
<div class="d-flex flex-column flex-md-row flex-items-center flex-md-items-center">
  <div class="col-2 d-flex flex-items-center flex-items-center flex-md-items-start">
    <img class="width-full avatar mb-2 mb-md-0" src="/github.png" alt="github" />
  </div>
  <div class="col-12 col-md-10 d-flex flex-column flex-justify-center flex-items-center flex-md-items-start pl-md-4">
    <h1 class="text-normal lh-condensed">GitHub</h1>
    <p class="h4 text-gray text-normal mb-2">How people build software.</p>
    <a class="text-gray text-small" href="#url">https://github.com/about</a>
  </div>
</div>
```


## Responsive grids
All the column width classes can be set per breakpoint to create responsive grid layouts. Each responsive style is applied to the specified breakpoint and up.

### Breakpoints
We use abbreviations for each breakpoint to keep the class names concise.

| Shorthand | Description |
| --- | --- |
| sm | min-width: 544px |
| md | min-width: 768px |
| lg | min-width: 1004px |
| xl | min-width: 1280px |

**Note:** The `lg` breakpoint matches our current page width of `980px` including left and right padding of `12px`. This is so that content doesn't touch the edges of the window when resized.

<hr />

In this example at the `sm` breakpoint 2 columns will show, at the `md` breakpoint 4 columns will show, and at the `lg` breakpoint 6 columns will show.

```html
<div class="container-lg clearfix">
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
  <div class="col-sm-6 col-md-3 col-lg-2 float-left p-2 border">
    .col-sm-6 .col-md-3 .col-lg-2
  </div>
</div>
```

For demonstration, this is how the above example would look at the `sm` breakpoint.

```html
<div class="container-lg clearfix">
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
  <div class="col-sm-6 float-left p-2 border">
    .col-sm-6
  </div>
</div>
```
This is how that same example would look at the `md` breakpoint.

```html
<div class="container-lg clearfix">
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
  <div class="col-md-3 float-left p-2 border">
    .col-md-3
  </div>
</div>
```

This is how that example would look at the `lg` breakpoint.

```html
<div class="container-lg clearfix">
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
  <div class="col-lg-2 float-left p-2 border">
    .col-lg-2
  </div>
</div>
```

## Containers
Container widths match our breakpoints and are available at a `md`, `lg`, and `xl` size. Containers apply a max-width rather than a fixed width for responsive layouts, and they center the container.

```html
<div class="container-md border">
  .container-md, max-width 768px
</div>

<div class="container-lg border">
  .container-lg, max-width 1012px
</div>

<div class="container-xl border">
  .container-xl, max-width 1280px
</div>
```

**Note:** `.container` is being replaced with `.container-lg`. To match the current fixed page width use `.container-lg` with `px-3`. This gives the container padding on smaller screens which works better for responsive layouts.
