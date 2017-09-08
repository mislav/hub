---
title: Filters
status: New release
status_issue: https://github.com/github/design-systems/issues/302
---

Filter utility classes can be applied to divs or images to apply visual effects.

<div class="flash flash-warn">
  Note: CSS filters are <a href="http://caniuse.com/#feat=css-filters">not supported by IE</a>
</div>

## Grayscale

Applying `.grayscale` to an element will remove all of its colors, and make it render in black and white.

```html
<img src="<%= image_path "modules/site/home-illo-business.svg" %>" class="img-responsive grayscale" alt="">
```
