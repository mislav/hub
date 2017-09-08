---
title: Spacing
status: Stable
source: https://github.com/primer/primer-css/blob/master/modules/primer-support/lib/variables/layout.scss
---

{:toc}

## Spacing scale
The spacing scale is a **base-8** scale. We chose a base-8 scale because eight is a highly composable number (it can be divided and multiplied many times and result in whole numbers), yet allows spacing dense enough for GitHub's UI. The scale's exception is that it begins at 4px to allow smaller padding and margin for denser parts of the site, from there on it steps up consistently in equal values of `8px`.

| Scale | Value |
| --- | --- |
| 0 | 0 |
| 1 | 4px |
| 2 | 8px |
| 3 | 16px |
| 4 | 24px |
| 5 | 32px |
| 6 | 40px |

The spacing scale is used for [margin](./utilities/margin) and [padding](./utilities/padding) utilities, and via variables within components.

## Em-based spacing
Ems are used for spacing within components such as buttons and form elements. We stick to common fractions for em values so that, in combination with typography and line-height, the total height lands on sensible numbers.

We aim for whole numbers, however, GitHub's body font-size is 14px which is difficult to work with, so we sometimes can't achieve a whole number. Less desirable values are highlighted in <span class="text-red">red</span> below.

| Fraction | Y Padding (em) | Total height at 14px | Total height at 16px |
| --- | --- | --- | --- |
| 3/4 | .75 | 42 | 48 |
| 1/2 | .5 | 35 | 40 |
| 3/8 | .375 | <span class="text-red">31.5</span> | 36 |
| 1/4 | .25 | 28 | 32 |
| 1/8 | .125 | <span class="text-red">24.5</span> | 28 |

We recommend using the fractions shown above. To calculate values with other font-sizes or em values, we suggest using [Formula](http://jxnblk.com/formula/).

## Spacer Variables

These variables match the above scale and are encouraged to be used within components. They are also used in our [margin](./utilities/margin) and [padding utilities](./utilities/padding).

```scss
$spacer-1: 4px;
$spacer-2: 8px;
$spacer-3: 16px;
$spacer-4: 24px;
$spacer-5: 32px;
$spacer-6: 40px;
```
