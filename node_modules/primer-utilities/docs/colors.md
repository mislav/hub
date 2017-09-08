---
title: Colors
status: New release
status_issue: https://github.com/github/design-systems/issues/97
---

Use color utilities to apply color to the background of elements, text, and borders.

* [Background colors](#background-colors)
* [Text colors](#text-colors)
* [Link colors](#link-colors)
* [Border colors](#border-colors)

## Background colors

Background colors are most commonly used for filling large blocks of content or areas with a color. When selecting a background color, make sure the foreground color contrast passes a minimum WCAG accessibility rating of [level AA](https://www.w3.org/TR/UNDERSTANDING-WCAG20/visual-audio-contrast-contrast.html). Meeting these standards ensures that content is accessible by everyone, regardless of disability or user device. You can [check your color combination with this demo site](http://jxnblk.com/colorable/demos/text/). For more information, read our [accessibility standards](/styleguide/css/principles/accessibility).

### Gray

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-gray</div>
    <code>#f5f5f5, $bg-gray</code>
    <div class="mt-2 bg-gray" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="col-6 float-left">
        <div class="h4">.bg-gray-dark</div>
        <code>#333, $bg-gray-dark</code>
        <div class="mt-2 bg-gray-dark border-right-0" style="height: 60px;"></div>
      </div>
      <div class="col-6 float-left">
        <div class="h4">.bg-gray-light</div>
        <code>#fafafa, $bg-gray-light</code>
        <div class="mt-2 bg-gray-light" style="height: 60px;"></div>
      </div>
    </div>
  </div>
</div>

### Blue

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-blue</div>
    <code>#4078c0, $bg-blue</code>
    <div class="mt-2 bg-blue" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="h4">.bg-blue-light</div>
      <code>#f2f8fa, $bg-blue-light</code>
      <div class="mt-2 bg-blue-light" style="height: 60px;"></div>
    </div>
  </div>
</div>

### Yellow

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-yellow</div>
    <code>#ffd36b, $bg-yellow</code>
    <div class="mt-2 bg-yellow" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="h4">.bg-yellow-light</div>
      <code>#fff9ea, $bg-yellow-light</code>
      <div class="mt-2 bg-yellow-light" style="height: 60px;"></div>
    </div>
  </div>
</div>

### Red

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-red</div>
    <code>#bd2c00, $bg-red</code>
    <div class="mt-2 bg-red" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="h4">.bg-red-light</div>
      <code>#fcdede, $bg-red-light</code>
      <div class="mt-2 bg-red-light" style="height: 60px;"></div>
    </div>
  </div>
</div>

### Green

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-green</div>
    <code>#6cc644, $bg-green</code>
    <div class="mt-2 bg-green" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="h4">.bg-green-light</div>
      <code>#eaffea, $bg-green-light</code>
      <div class="mt-2 bg-green-light" style="height: 60px;"></div>
    </div>
  </div>
</div>


### Purple

<div class="container-lg clearfix mb-4">
  <div class="col-3 float-left pr-4">
    <div class="h4">.bg-purple</div>
    <code>#6e5494, $bg-purple</code>
    <div class="mt-2 bg-purple" style="height: 60px;"></div>
  </div>
  <div class="col-9 float-left">
    <div class="container-lg clearfix">
      <div class="h4">.bg-purple-light</div>
      <code>#f5f0ff, $bg-purple-light</code>
      <div class="mt-2 bg-purple-light" style="height: 60px;"></div>
    </div>
  </div>
</div>

## Text colors

Use text color utilities to set text or [octicons](/styleguide/css/styles/core/components/octicons) to a specific color. Color contrast must pass a minimum WCAG accessibility rating of [level AA](https://www.w3.org/TR/UNDERSTANDING-WCAG20/visual-audio-contrast-contrast.html). This ensures that viewers who cannot see the full color spectrum are able to read the text. To customize outside of the suggested combinations below, we recommend using this [color contrast testing tool](http://jxnblk.com/colorable/demos/text/). For more information, read our [accessibility standards](/styleguide/css/principles/accessibility).

These are our most common text with background color combinations. They don't all pass accessibility standards currently, but will be updated in the future. **Any of the combinations with a warning icon must be used with caution**.

### Text on white background

```html
<div class="text-blue mb-2">
  .text-blue on white
</div>
<div class="text-gray-dark mb-2">
  .text-gray-dark on white
</div>
<div class="text-gray mb-2">
  .text-gray on white
</div>
<div class="text-red mb-2">
  .text-red on white
</div>
<div class="text-orange mb-2">
  .text-orange on white
</div>
<span class="float-left text-red tooltipped tooltipped-n" aria-label="Does not meet accessibility standards"><%= octicon("alert") %></span>
<div class="text-orange-light mb-2">
  .text-orange-light on white
</div>
<span class="float-left text-red tooltipped tooltipped-n" aria-label="Does not meet accessibility standards"><%= octicon("alert") %></span>
<div class="text-green mb-2 ml-4">
  .text-green on white
</div>
<div class="text-purple mb-2">
  .text-purple on white
</div>
```

### Text on colored backgrounds

```html
<div class="text-white bg-blue mb-2">
  .text-white on .bg-blue
</div>
<div class="bg-blue-light mb-2">
  .text-gray-dark on .bg-blue-light
</div>
<div class="text-white bg-red mb-2">
  .text-white on .bg-red
</div>
<div class="text-red bg-red-light mb-2">
  .text-red on .bg-red-light
</div>
<div class="bg-green-light mb-2">
  .text-gray-dark on .bg-green-light
</div>
<div class="bg-yellow mb-2">
  .text-gray-dark on .bg-yellow
</div>
<div class="bg-yellow-light mb-2">
  .text-gray-dark on .bg-yellow-light
</div>
<div class="text-white bg-purple mb-2">
  .text-white on .bg-purple
</div>
<div class="text-white bg-gray-dark mb-2">
  .text-white on .bg-gray-dark
</div>
<div class="bg-gray">
  .text-gray-dark on .bg-gray
</div>
```

## Link colors

Base link styles turn links blue and apply an underline on hover. You can override the base link styles with text color utilities and the following link utilities. **Bear in mind that link styles are easier for more people to see and interact with when the changes in styles do not rely on color alone.**

Use `link-gray` to turn the link color to `$text-gray` and remain hover on blue.

```html
<a class="link-gray" href="#url">link-gray</a>
```

Use `link-gray-dark` to turn the link color to `$text-gray-dark` and remain hover on blue.

```html
<a class="link-gray-dark"  href="#url">link-gray-dark</a>
```

Use `.muted-link` to turn the link light gray in color, and blue on hover or focus with no underline.

```html
<a class="muted-link" href="#url">muted-link</a>
```

Use `link-hover-blue` to make any text color used with links to turn blue on hover. This is useful when you want only part of a link to turn blue on hover.

```html
<a class="text-gray-dark no-underline" href="#url">
  A link with only part of it is <span class="link-hover-blue">blue on hover</span>.
</a>
```

## Border colors

Border colors are documented on the [border utilities page](/styleguide/css/styles/core/utilities/borders#border-width-style-and-color-utilities).
