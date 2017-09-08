---
title: Typography
status_issue: https://github.com/github/design-systems/issues/329
status: New release
source: https://github.com/primer/primer-css/blob/master/modules/primer-support/lib/variables/typography.scss
---

{:toc}

## Type Scale

The typography scale is designed to work for GitHub's product UI and marketing sites. Font sizes are designed to work in combination with line-height values so as to result in more sensible numbers wherever possible.

Font sizes are smaller on mobile and scale up at the `md` [breakpoint](#breakpoints) to be larger on desktop.

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

The typography scale is used to create [typography utilities](./utilities/typography).

## Typography variables

#### Font size variables
```scss
// Heading sizes - mobile
// h4—h6 remain the same size on both mobile & desktop
$h00-size-mobile: 40px;
$h0-size-mobile: 32px;
$h1-size-mobile: 26px;
$h2-size-mobile: 22px;
$h3-size-mobile: 18px;

// Heading sizes - desktop
$h00-size: 48px;
$h0-size: 40px;
$h1-size: 32px;
$h2-size: 24px;
$h3-size: 20px;
$h4-size: 16px;
$h5-size: 14px;
$h6-size: 12px;
```

#### Font weight variables
```scss
$font-weight-bold: 600 !default;
$font-weight-light: 300 !default;
```

#### Line height variables
```scss
$lh-condensed-ultra: 1 !default;
$lh-condensed: 1.25 !default;
$lh-default: 1.5 !default;
```

## Typography Mixins
Typography mixins are available for heading styles and for our type scale. They can be used within components or custom CSS. The same styles are also available as [utilities](./utilities/typography#heading-utilities) which requires no additional CSS.

Heading mixins are available for `h1` through to `h6`, this includes the font-size and font-weight. Example:

```scss
.styles {
  @include h1;
}
```

Responsive heading mixins can be used to adjust the font-size between the `sm` and `lg` breakpoint. Only headings 1—3 are responsive. Heading 4—6 are small enough to remain the same between mobile and desktop. Responsive heading mixins include the font-size and font-weight as well as the media query. To use a responsive heading mixin, postfix the heading with `-responsive`:

```scss
.styles {
  @include h1-responsive;
}
```

Responsive type scale mixins are also available. Type scale mixins apply a font-size that gets slightly bigger at the `lg` breakpoint. They do not apply a font-weight. Like the responsive heading mixins, they are only available for size `f1` to `f3` and are postfixed with `-responsive` as in the example below:

```scss
.style {
  @include f1-responsive;
}
```
