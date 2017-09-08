---
title: Typography
status: Stable
---

Type utilities are designed to work in combination with line-height utilities so as to result in more sensible numbers wherever possible. These also exist as [variables](/styleguide/css/styles/core/support/typography#typography-variables) that you can use in components or custom CSS.

{:toc}

Font sizes are smaller on mobile and scale up at the `md` [breakpoint](/styleguide/css/styles/core/support/breakpoints) to be larger on desktop.

| Scale | Font size: mobile | Font size: desktop | 1.25 line height | 1.5 line height |
| --- | --- | --- | --- | --- |
| 00 | 40px | 48px | 60 | 72 |
| 0 | 32px | 40px | 50 | 60 |
| 1 | 26px | 32px | 40 | 48 |
| 2 | 22px | 24px | 30 | 36 |
| 3 | 18px | 20px | 25 | 30 |
| 4 | 16px | 16px | 20 | 24 |
| 5 | 14px | 14px | 17.5 | 21 |
| 6 | 12px | 12px | 15 | 18 |


## Heading utilities

Use `.h1` – `.h6` to change an elements font size and weight to match our heading styles.

```html
<p class="h1">Pizza 1</p>
<p class="h2">Pizza 2</p>
<p class="h3">Pizza 3</p>
<p class="h4">Pizza 4</p>
<p class="h5">Pizza 5</p>
<p class="h6">Pizza 6</p>
```

These are particularly useful for changing the visual appearance while keeping the markup semantically correct. Be sure you keep the hierarchy appropriate for the page.

```html
<h2 class="h1">Looks like a heading 1, semantically a heading 2</h2>
```

## Type scale utilities

Use `.f1` – `.f6` to change an elements font size while keeping inline with our type scale.

```html
<p class="f1">Focaccia</p>
<p class="f2">Focaccia</p>
<p class="f3">Focaccia</p>
<p class="f4">Focaccia</p>
<p class="f5">Focaccia</p>
<p class="f6">Focaccia</p>
```

Lighter font-weight utilities are available in a limited range. Lighter font-weights reduce the legibility of text, particularly at small font sizes, so the scale only goes down to `f3` at 20px. The larger sizes`f0` and `f00` allow for lighter and larger type that is in keeping with our marketing styles.

```html
<h1 class="f00-light">Potato chips</h1>
<h1 class="f0-light">Potato chips</h1>
<h1 class="f1-light">Potato chips</h1>
<h1 class="f2-light">Potato chips</h1>
<h1 class="f3-light">Potato chips</h1>
```

## Line height styles
Change the line height density with these utilities.

```html
<p class="lh-default">
  Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail. Jowl chuck biltong flank meatball, beef short ribs. Jowl tenderloin ground round, short loin tri-tip ribeye picanha filet mignon pig chicken kielbasa t-bone fatback. Beef ribs meatball chicken corned beef salami.
</p>
<p class="lh-condensed">
  Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail. Jowl chuck biltong flank meatball, beef short ribs. Jowl tenderloin ground round, short loin tri-tip ribeye picanha filet mignon pig chicken kielbasa t-bone fatback. Beef ribs meatball chicken corned beef salami.
</p>
<p class="lh-condensed-ultra">
  Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail. Jowl chuck biltong flank meatball, beef short ribs. Jowl tenderloin ground round, short loin tri-tip ribeye picanha filet mignon pig chicken kielbasa t-bone fatback. Beef ribs meatball chicken corned beef salami.
</p>
```

## Typographic styles
Change the font weight, styles, and alignment with these utilities.

```html
<p class="text-normal">Normal</p>
<p class="text-italic">Italic</p>
<p class="text-bold">Bold</p>
<p class="text-uppercase">Uppercase</p>
<p class="no-wrap">No wrap</p>
<p class="ws-normal">Normal whitespace</p>
<p class="wb-break-all">Line break long lines</p>
<p class="no-underline">No underline</p>
<p class="text-emphasized">Emphasized</p>
<p class="text-small">Small</p>
<p class="lead">Bacon ipsum dolor amet tri-tip chicken kielbasa, cow swine beef corned beef ground round prosciutto hamburger porchetta sausage alcatra tail.</p>
```

## Text alignment

 Use text alignment utilities to left align, center, or right align text.

```html
<p class="text-left">Left align</p>
<p class="text-center">Center</p>
<p class="text-right">Right align</p>
```
## Responsive text alignment

Use the following formula to make a text alignment utility responsive: `.text-[breakpoint]-[alignment]`

```html
<p class="text-left text-sm-center text-lg-right">Left in a small viewport, centered in mid-sized viewports, and right aligned in larger viewports</p>
```

## List styles

Remove bullets from an unordered list or numbers from an ordered list by applying `.list-style-none` to the `<ul>`.

```html
<ul class="list-style-none">
  <li>First list item</li>
  <li>Second list item</li>
  <li>Third list item</li>
</ul>
```

## Text Shadows

Text shadows can be used to help readability and to add some depth on colored backgrounds.

```html
<div class="bg-gray-dark p-5">
  <h3 class="text-white text-shadow-dark">.text-shadow-dark helps white text stand out on dark or photographic backgrounds</h3>
</div>
<div class="bg-gray p-5">
  <h3 class="text-shadow-light">.text-shadow-light creates an embossed effect for dark text</h3>
</div>
```
