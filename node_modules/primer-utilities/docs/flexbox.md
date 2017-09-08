---
title: Flexbox
status: New release
source: /app/assets/stylesheets/primer-core/utilities/lib/layout.scss
status_issue: https://github.com/github/design-systems/issues/157
---

Flex utilities can be used to apply `flexbox` behaviors to elements by using utility classes.

{:toc}

## Required reading

Before using these utilities, **you should be familiar with CSS3 Flexible Box spec**. If you are not, check out MDN's guide:

:book: **[Using CSS Flexible Boxes](https://developer.mozilla.org/en-US/docs/Web/CSS/CSS_Flexible_Box_Layout/Using_CSS_flexible_boxes)**

## Flex container

Use these classes to make an element lay out its content using the flexbox model. Each **direct** child of the flex container will become a flex item.

#### CSS

```css
.d-flex { display: flex; }
.d-inline-flex { display: inline-flex; }
```

#### Classes

| Class | Description |
| --- | --- |
| `.d-flex` | The element behaves like a block and lays out its content using the flexbox model. |
| `.d-inline-flex` | The element behaves like an inline element and lays out its content using the flexbox model. |

#### Example using `.d-flex`

```html
<!-- flex container -->
<div class="border d-flex">
  <div class="p-5 border bg-gray-light">flex item 1</div>
  <div class="p-5 border bg-gray-light">flex item 2</div>
  <div class="p-5 border bg-gray-light">flex item 3</div>
</div>
```

#### Example using `.d-inline-flex`

```html
<!-- inline-flex container -->
<div class="border d-inline-flex">
  <div class="p-5 border bg-gray-light">flex item 1</div>
  <div class="p-5 border bg-gray-light">flex item 2</div>
  <div class="p-5 border bg-gray-light">flex item 3</div>
</div>
```

## Flex direction

Use these classes to define the orientation of the main axis (`row` or `column`). By default, flex items will display in a row. Use `.flex-column` to change the direction of the main axis to vertical.

#### CSS

```css
.flex-row         { flex-direction: row; }
.flex-row-reverse { flex-direction: row-reverse; }
.flex-column      { flex-direction: column; }
```

#### Classes

| Class | Description |
| --- | --- |
| `.flex-row` | The main axis runs left to right (default). |
| `.flex-row-reverse` | The main axis runs right to left. |
| `.flex-column` | The main axis runs top to bottom. |

#### Example using `.flex-column`

```html
<div class="border d-flex flex-column">
  <div class="p-5 border bg-gray-light">Item 1</div>
  <div class="p-5 border bg-gray-light">Item 2</div>
  <div class="p-5 border bg-gray-light">Item 3</div>
</div>
```

#### Example using `.flex-row`

This example uses the responsive variant `.flex-md-row` to override `.flex-column` Learn more about responsive flexbox utilities **[here](#responsive-flex-utilities)**.

```html
<div class="border d-flex flex-column flex-md-row">
  <div class="p-5 border bg-gray-light">Item 1</div>
  <div class="p-5 border bg-gray-light">Item 2</div>
  <div class="p-5 border bg-gray-light">Item 3</div>
</div>
```

#### Example using `.flex-row-reverse`

This example uses the responsive variant `.flex-md-row-reverse` to override `.flex-column` Learn more about responsive flexbox utilities **[here](#responsive-flex-utilities)**.

```html
<div class="border d-flex flex-column flex-md-row-reverse">
  <div class="p-5 border bg-gray-light">Item 1</div>
  <div class="p-5 border bg-gray-light">Item 2</div>
  <div class="p-5 border bg-gray-light">Item 3</div>
</div>
```

## Flex wrap

You can choose whether flex items are forced into a single line or wrapped onto multiple lines.

#### CSS

```css
.flex-wrap     { flex-wrap: wrap; }
.flex-nowrap   { flex-wrap: nowrap; }
```

#### Classes

| Class | Description |
| --- | --- |
| `.flex-wrap` | Flex items will break onto multiple lines (default) |
| `.flex-nowrap` | Flex items are laid out in a single line, even if they overflow |

#### `flex-wrap` example

```html
<div class="border d-flex flex-wrap">
  <div class="p-5 px-6 border bg-gray-light">1</div>
  <div class="p-5 px-6 border bg-gray-light">2</div>
  <div class="p-5 px-6 border bg-gray-light">3</div>
  <div class="p-5 px-6 border bg-gray-light">4</div>
  <div class="p-5 px-6 border bg-gray-light">5</div>
  <div class="p-5 px-6 border bg-gray-light">6</div>
  <div class="p-5 px-6 border bg-gray-light">7</div>
  <div class="p-5 px-6 border bg-gray-light">8</div>
  <div class="p-5 px-6 border bg-gray-light">9</div>
</div>
```

#### `flex-nowrap` example

```html
<div class="border d-flex flex-nowrap">
  <div class="p-5 px-6 border bg-gray-light">1</div>
  <div class="p-5 px-6 border bg-gray-light">2</div>
  <div class="p-5 px-6 border bg-gray-light">3</div>
  <div class="p-5 px-6 border bg-gray-light">4</div>
  <div class="p-5 px-6 border bg-gray-light">5</div>
  <div class="p-5 px-6 border bg-gray-light">6</div>
  <div class="p-5 px-6 border bg-gray-light">7</div>
  <div class="p-5 px-6 border bg-gray-light">8</div>
  <div class="p-5 px-6 border bg-gray-light">9</div>
</div>
```

## Justify content

Use these classes to distribute space between and around flex items along the **main axis** of the container.

#### CSS

```CSS
.flex-justify-start    { justify-content: flex-start; }
.flex-justify-end      { justify-content: flex-end; }
.flex-justify-center   { justify-content: center; }
.flex-justify-between  { justify-content: space-between; }
.flex-justify-around   { justify-content: space-around; }
```

#### Classes

| Class | Default behavior |
| --- | --- |
| `.flex-justify-start` | Justify all items to the left |
| `.flex-justify-end` | Justify all items to the right |
| `.flex-justify-center` | Justify items to the center of the container |
| `.flex-justify-between` | Distribute items evenly. First item is on the start line, last item is on the end line. |
| `.flex-justify-around` | Distribute items evenly with equal space around them |

#### flex-justify-start

Use `.flex-justify-start` to align items to the start line. By default this will be on the left side of the container. If you are using `.flex-column`, the end line will be at the top of the container.

```html
<div class="border d-flex flex-justify-start">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
</div>
```

#### flex-justify-end

Use `.flex-justify-end` to align items to the end line. By default this will be on the right side of the container. If you are using `.flex-column`, the end line will be at the bottom of the container.

```html
<div class="border d-flex flex-justify-end">
 <div class="p-5 border bg-gray-light">1</div>
 <div class="p-5 border bg-gray-light">2</div>
 <div class="p-5 border bg-gray-light">3</div>
</div>
```

#### flex-justify-center

Use `.flex-justify-center` to align items in the middle of the container.

```html
<div class="border d-flex flex-justify-center">
 <div class="p-5 border bg-gray-light">1</div>
 <div class="p-5 border bg-gray-light">2</div>
 <div class="p-5 border bg-gray-light">3</div>
</div>
```

#### flex-justify-between

Use `.flex-justify-between` to distribute items evenly in the container. The first items will be on the start line and the last item will be on the end line.

```html
<div class="border d-flex flex-justify-between">
 <div class="p-5 border bg-gray-light">1</div>
 <div class="p-5 border bg-gray-light">2</div>
 <div class="p-5 border bg-gray-light">3</div>
</div>
```

#### flex-justify-around

Use `.flex-justify-around` to distribute items evenly, with an equal amount of space around them. Visually the spaces won't look equal, since each item has the same unit of space around it. For example, the first item only has one unit of space between it and the start line, but it has two units of space between it and the next item.

```html
<div class="border d-flex flex-justify-around">
 <div class="p-5 border bg-gray-light">1</div>
 <div class="p-5 border bg-gray-light">2</div>
 <div class="p-5 border bg-gray-light">3</div>
</div>
```

## Align items

Use these classes to align items on the **cross axis**.

The cross axis runs perpendicular to the main axis. By default the main axis runs horizontally, so your items will align vertically on the cross axis. If you use [flex-column](#flex-direction) to make the main axis run vertically, your items will be aligned horizontally.

#### CSS

```css
.flex-items-start      { align-items: flex-start; }
.flex-items-end        { align-items: flex-end; }
.flex-items-center     { align-items: center; }
.flex-items-baseline   { align-items: baseline; }
.flex-items-stretch    { align-items: stretch; }
```

#### Classes

| Class | Behavior |
| --- | --- |
| `.flex-items-start` | Align items to the start of the cross axis  |
| `.flex-items-end` | Align items to the end of the cross axis |
| `.flex-items-center` | Align items to the center of the cross axis |
| `.flex-items-baseline` | Align items along their baselines |
| `.flex-items-stretch` | Stretch items from start of cross axis to end of cross axis |

#### flex-items-start

```html
<div  style="min-height: 150px;" class="border d-flex flex-items-start">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
</div>
```

#### flex-items-end

```html
<div  style="min-height: 150px;" class="border d-flex flex-items-end">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
</div>
```

#### flex-items-center

```html
<div  style="min-height: 150px;" class="border d-flex flex-items-center">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
</div>
```

#### flex-items-baseline

With `.flex-items-baseline`, items will be aligned along their baselines even if they have different font sizes.

```html
<div  style="min-height: 150px;" class="border d-flex flex-items-baseline">
  <div class="p-5 border bg-gray-light f1">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light f3">3</div>
  <div class="p-5 border bg-gray-light">4</div>
</div>
```

#### flex-items-stretch

```html
<div  style="min-height: 150px;" class="border d-flex flex-items-stretch">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
</div>
```

## Align content

When the main axis wraps, this creates multiple main axis lines and adds extra space in the cross axis. Use these classes to align the lines of the main axis on the cross axis.

#### CSS

```css
.flex-content-start    { align-content: flex-start; }
.flex-content-end      { align-content: flex-end; }
.flex-content-center   { align-content: center; }
.flex-content-between  { align-content: space-between; }
.flex-content-around   { align-content: space-around; }
.flex-content-stretch  { align-content: stretch; }
```

#### Classes

| Class | Description |
| --- | --- |
| `.flex-content-start` | Align content to the start of the cross axis  |
| `.flex-content-end` | Align content to the end of the cross axis |
| `.flex-content-center` | Align content to the center of the cross axis |
| `.flex-content-between` | Distribute content evenly. First line is on the start of the cross axis, last line is on the end of the cross axis.  |
| `.flex-content-around` | Stretch items from the start of the cross axis to the end of the cross axis. |
| `.flex-content-stretch` | Lines stretch to occupy available space.  |

#### flex-content-start

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-start">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

#### flex-content-end

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-end">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

#### flex-content-center

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-center">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

#### flex-content-between

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-between">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

#### flex-content-around

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-around">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

#### flex-content-stretch

```html
<div style="min-height: 300px;" class="border d-flex flex-wrap flex-content-stretch">
  <div class="p-5 border bg-gray-light">1</div>
  <div class="p-5 border bg-gray-light">2</div>
  <div class="p-5 border bg-gray-light">3</div>
  <div class="p-5 border bg-gray-light">4</div>
  <div class="p-5 border bg-gray-light">5</div>
  <div class="p-5 border bg-gray-light">6</div>
  <div class="p-5 border bg-gray-light">7</div>
  <div class="p-5 border bg-gray-light">8</div>
  <div class="p-5 border bg-gray-light">9</div>
  <div class="p-5 border bg-gray-light">10</div>
  <div class="p-5 border bg-gray-light">11</div>
  <div class="p-5 border bg-gray-light">12</div>
</div>
```

## Flex

Use this class to specify the ability of a flex item to alter its dimensions to fill available space.

```CSS
.flex-auto    { flex: 1 1 auto; }
```

| Class | Description |
| --- | --- |
| `.flex-auto` | Sets default values for  `flex-grow` (1), `flex-shrink` (1) and `flex-basis` (auto)  |

#### flex-auto

```html
<div class="border d-flex">
  <div class="p-5 border bg-gray-light flex-auto">.flex-auto</div>
  <div class="p-5 border bg-gray-light flex-auto">.flex-auto</div>
  <div class="p-5 border bg-gray-light flex-auto">.flex-auto</div>
</div>
```

## Align self

Use these classes to adjust the alignment of an individual flex item on the cross axis. This overrides any `align-items` property applied to the flex container.

#### CSS

```css
.flex-self-auto        { align-self: auto; }
.flex-self-start       { align-self: flex-start; }
.flex-self-end         { align-self: flex-end; }
.flex-self-center      { align-self: center; }
.flex-self-baseline    { align-self: baseline; }
.flex-self-stretch     { align-self: stretch; }
```

#### Classes

| Class | Description |
| --- | --- |
| `.flex-self-auto` | Inherit alignment from parent |
| `.flex-self-start` | Align to the start of the cross axis  |
| `.flex-self-end` | Align to the end of the cross axis |
| `.flex-self-center` | Align to center of cross axis |
| `.flex-self-baseline` | Align baseline to the start of the cross axis |
| `.flex-self-stretch` | Stretch item from start of cross axis to end of cross axis. |

#### flex-self-auto

```html
<div style="min-height: 150px;" class="border d-flex">
  <div class="p-5 border bg-gray-light flex-self-auto">.flex-self-auto</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

#### flex-self-start

```html
<div style="min-height: 150px;" class="border d-flex">
  <div class="p-5 border bg-gray-light flex-self-start">.flex-self-start</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

#### flex-self-end

```html
<div style="min-height: 150px;" class="border d-flex">
  <div class="p-5 border bg-gray-light flex-self-end">.flex-self-end</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

#### flex-self-center

```html
<div style="min-height: 150px;" class="border d-flex">
  <div class="p-5 border bg-gray-light flex-self-center">.flex-self-center</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

#### flex-self-baseline

```html
<div style="min-height: 150px;" class="border d-flex flex-items-end">
  <div class="p-5 border bg-gray-light flex-self-baseline f4">.flex-self-baseline</div>
  <div class="p-5 border bg-gray-light f1">item</div>
  <div class="p-5 border bg-gray-light f00">item</div>
</div>
```

#### flex-self-stretch

```html
<div style="min-height: 150px;" class="border d-flex flex-items-start">
  <div class="p-5 border bg-gray-light flex-self-stretch">.flex-self-stretch</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

## Responsive flex utilities

All flexbox utilities can be adjust per [breakpoint](/styleguide/css/modules/grid#breakpoints) using the following formulas:

- `d-[breakpoint]-[property]` for `display`
- `flex-[breakpoint]-[property]-[behavior]` for various flex properties
- `flex-[breakpoint]-item-equal` for equal width and equal height flex items

Each responsive style is applied to the specified breakpoint and up.

#### Example classes

```css
/* Example classes */

.d-sm-flex {}
.d-md-inline-flex {}

.flex-lg-row {}
.flex-xl-column {}

.flex-sm-wrap {}
.flex-lg-nowrap {}

.flex-lg-self-start {}

.flex-md-item-equal {}
```

#### Example usage

Mixing flex properties:

```html
<div style="min-height: 150px;" class="border d-flex flex-items-start flex-md-items-center flex-justify-start flex-lg-justify-between">
  <div class="p-5 border bg-gray-light flex-md-self-stretch">.flex-self-stretch</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
  <div class="p-5 border bg-gray-light">&nbsp;</div>
</div>
```

Using the equal width, equal height utilities:

```html
<div class="border d-flex">
  <div class="flex-md-item-equal p-3 border bg-gray-light">Stuff and things</div>
  <div class="flex-md-item-equal p-3 border bg-gray-light">More stuff<br> on multiple lines</div>
  <div class="flex-md-item-equal p-3 border bg-gray-light">Hi mom</div>
</div>
```

## Example components

The flex utilities can be used to create a number of components. Here are some examples.

### Media object

You can use flex utilities to make a simple media object, like this:

```html
<div class="border d-flex flex-items-center p-4">
  <div class="p-5 border bg-gray-light">image</div>
  <p class="pl-4 m-0"><b>Body</b> Bacon ipsum dolor sit amet chuck prosciutto landjaeger ham hock filet mignon shoulder hamburger pig venison.</p>
</div>
```

### Responsive media object

Here is an example of a media object that is **vertically centered on large screens**, but converts to a stacked column layout on small screens.

```html
<div class="border p-3 d-flex flex-column flex-md-row flex-md-items-center">
  <div class="pr-0 pr-md-3 mb-3 mb-md-0 d-flex flex-justify-center flex-md-justify-start">
    <img style="max-width:100px; max-height:100px;" src="/images/gravatars/gravatar-user-420.png" />
  </div>
  <div class="d-flex text-center text-md-left">
    <p><b>Body</b> Bacon ipsum dolor sit amet chuck prosciutto landjaeger ham hock filet mignon shoulder hamburger pig venison.</p>
  </div>
  <div class="ml-md-3 d-flex flex-justify-center">
    <%= octicon "mark-github", :height => '32' %>
  </div>
</div>
```

## Flexbox bugs

This section lists flexbox bugs that affect browsers we [currently support](.../styles#user-content-browser-support).

**1. Minimum content sizing of flex items not honored:** Some browsers don't respect flex item size. Instead of respecting the minimum content size, items shrink below their minimum size which can create some undesirable results, such as overflowing text. The workaround is to apply `flex-shrink: 0;` to the items using `d-flex`. This can be applied with the `flex-shrink-0` utility. For more information read [philipwalton/flexbugs](https://github.com/philipwalton/flexbugs#1-minimum-content-sizing-of-flex-items-not-honored).
