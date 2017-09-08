---
title: Layout
status: Stable
---

Change the document layout with display, float, alignment, and other utilities.

{:toc}

## Display
Adjust the display of an element with `.d-block`, `.d-none`, `.d-inline`, `.d-inline-block`, `.d-table`, `.d-table-cell` utilities.

```html
<div class="d-inline border">
  .d-inline
  <div class="d-inline-block border">
    .d-inline-block
  </div>
  <span class="d-block border">.d-block</span>
  <div class="d-table border">
    <div class="d-table-cell border">
      .d-table-cell
    </div>
  </div>
  <div class="d-table-cell border">
    .d-table-cell
  </div>
  <div class="d-none">
    .d-none
  </div>
</div>
```

There are known issues with using `display:table` and wrapping long strings, particularly in Firefox. You may need to use `table-fixed` on elements with `d-table` and apply column widths to table cells, which you can do with our [column width styles](/styleguide/css/modules/grid#column-widths).

```html
<div class="d-table table-fixed width-full">
  <div class="d-table-cell border">
    .d-table-cell
  </div>
  <div class="d-table-cell col-10 border">
    d-table-cell .col-10
  </div>
</div>
```

### Responsive display
A selection of display utilities are able to be applied or changed per [breakpoint](/styleguide/css/modules/grid#breakpoints). `.d-block`, `.d-none`, `.d-inline`, and `.d-inline-block` are available as responsive utilities using the following formula: `d-[breakpoint]-[property]`. For example: `d-md-inline-block`. Each responsive display utility is applied to the specified breakpoint and up.

In the following example, the `ul` element switches from `display: block` on mobile to  `display: inline-block` at the `md` breakpoint, while the list items remain inline.

```html
<h5 class="d-md-inline-block">.d-md-inline-block</h5>
<ul class="d-md-inline-block border">
  <li class="d-inline border">.d-inline</li>
  <li class="d-inline border">.d-inline</li>
</ul>
```

### Responsive hide
Hide utilities are able to be applied or changed per breakpoint using the following formula:<br /> `hide-[breakpoint]`, for example: `hide-sm`. Hide utilities act differently from other responsive styles and are applied to each **breakpoint-range only**.

| Shorthand | Range |
| --- | --- |
| -sm | 0—544px |
| -md | 544px—768px |
| -lg | 768px—1004px |
| -xl | 1004px and above |

```html
<div>
  <h3>Potato update</h3>
  <span class="hide-sm hide-md">Opened by <a href="#url">broccolini</a></span>
  <span class="d-md-none">Updated</span> 3 hours ago
</div>
```

### Text direction
`.direction-ltr` or `.direction-rtl` can be used to change the text direction. This is especially helpful when paired with `.d-table`, `.d-table-cell`, and `.v-align-middle` to create equal height, vertically centered, alternating content.

## Visibility
Adjust the visibility of an element with `.v-hidden` and `.v-visible`.

## Overflow
Adjust element overflow with `.overflow-hidden`, `.overflow-scroll`, and `.overflow-auto`. `.overflow-hidden` can also be used to create a new [block formatting context](https://developer.mozilla.org/en-US/docs/Web/Guide/CSS/Block_formatting_context) or clear floats.

## Floats
Use `.float-left` and `.float-right` to set floats, and `.clearfix` to clear.
```html
<div class="clearfix border border-gray">
  <div class="float-left border border-gray">
    .float-left
  </div>
  <div class="float-right border border-gray">
    .float-right
  </div>
</div>
```
### Responsive floats
Float utilities can be applied or changed per [breakpoint](/styleguide/css/modules/grid#breakpoints). This can be useful for responsive layouts when you want an element to be full width on mobile but floated at a larger breakpoint.

Each responsive float utility is applied to the specified breakpoint and up, using the following formula:  `float-[breakpoint]-[property]`. For example: `float-md-left`. Remember to use `.clearfix` to clear.

```html
<div class="clearfix border border-gray">
  <div class="float-md-left border border-gray">
    .float-md-left
  </div>
  <div class="float-md-right border border-gray">
    .float-md-right
  </div>
</div>
```

## Alignment
Adjust the alignment of an element with `.v-align-baseline`, `.v-align-top`, `.v-align-middle` or `.v-align-bottom`. The vertical-align property only applies to inline or table-cell boxes.

```html
<div class="d-table border border-gray">
  <div class="d-table-cell"><h1>Potatoes</h1></div>
  <div class="d-table-cell v-align-baseline">.v-align-baseline</div>
  <div class="d-table-cell v-align-top">.v-align-top</div>
  <div class="d-table-cell v-align-middle">.v-align-middle</div>
  <div class="d-table-cell v-align-bottom">.v-align-bottom</div>
</div>
```

Use `v-align-text-top` or `v-align-text-bottom` to adjust the alignment of an element with the top or bottom of the parent element's font.

```html
<div class="border border-gray">
  <h1 class="mr-1">Potatoes
    <span class="f4 v-align-text-top mr-1">.v-align-text-top</span>
    <span class="f4 v-align-text-bottom mr-1">.v-align-text-bottom</span>
  </h1>
</div>
```

## Width and height

Use `.width-fit` to set max-width 100%.

```html
<div class="one-fourth column">
  <img class="width-fit bg-gray" src="/images/gravatars/gravatar-user-420.png" alt="width fitted octocat" />
</div>
```

Use `.width-full` to set width to 100%.

```html
<div class="d-table width-full">
  <div class="d-table-cell">
    <input class="form-control width-full" type="text" value="Full width form field" aria-label="Sample full width form field">
  </div>
</div>
```

Use `.height-full` to set height to 100%.

```html
<div class="d-table border border-gray">
  <div class="d-table-cell height-full v-align-middle pl-3">
    <%= octicon "three-bars" %>
  </div>
  <div class="p-3">
    Bacon ipsum dolor amet meatball flank beef tail pig boudin ham hock chicken capicola. Shoulder ham spare ribs turducken pork tongue. Bresaola corned beef sausage jowl ribeye kielbasa tenderloin andouille leberkas tongue. Ribeye tri-tip tenderloin pig, chuck ground round chicken tongue corned beef biltong.
  </div>
</div>
```

## Position
Position utilities can be used to alter the default document flow. **Be careful when using positioning, it's often unnecessary and commonly misused.**

Use `.position-relative` to create a new stacking context.

_Note how the other elements are displayed as if "Two" were in its normal position and taking up space._
```html
<div class="d-inline-block float-left bg-blue text-white m-3" style="width:100px; height:100px;">
  One
</div>
<div class="d-inline-block float-left position-relative bg-blue text-white m-3" style="width:100px; height:100px;  top:12px; left:12px;">
  Two
</div>
<div class="d-inline-block float-left bg-blue text-white m-3" style="width:100px; height:100px;">
  Three
</div>
<div class="d-inline-block float-left bg-blue text-white m-3" style="width:100px; height:100px;">
  Four
</div>
```

Use `.position-absolute` to take elements out of the normal document flow.

```html
<div class="position-relative" style="height:116px;">
  <button type="button" class="btn btn-secondary mb-1">Button</button>
  <div class="position-absolute border border-gray p-2">
    <a href="#url" class="d-block p-1">Mashed potatoes</a>
    <a href="#url" class="d-block p-1">Fries</a>
  </div>
</div>
```

Use `.position-fixed` to position an element relative to the viewport. **Be careful when using fixed positioning. It is tricky to use and can lead to unwanted side effects.**

_Note: fixed positioning has been disabled here for demonstration only._

```html
<div class="position-fixed bg-gray-light border-bottom border-gray p-3">
  .position-fixed
</div>
```

Use `top-0`, `right-0`, `bottom-0`, and `left-0` with `position-absolute`, `position-fixed`, or `position-relative` to further specify an elements position.

```html
<div style="height: 64px;">
  <div class="border position-absolute top-0 left-0">
    .top-0 .left-0
  </div>
  <div class="border position-absolute top-0 right-0">
    .top-0 .right-0
  </div>
  <div class="border position-absolute bottom-0 right-0">
    .bottom-0 .right-0
  </div>
  <div class="border position-absolute bottom-0 left-0">
    .bottom-0 .left-0
  </div>
</div>
```

To fill an entire width or height, use opposing directions.

_Note: fixed positioning has been disabled here for demonstration only._

```html
<div class="position-fixed left-0 right-0 p-3 bg-gray-dark text-white">
  .position-fixed .left-0 .right-0
</div>
```

### Screen reader only

Use `.sr-only` to position an element outside of the viewport for screen reader access only. **Even though the element can't be seen, make sure it still has a sensible tab order.**

```html
<div class="js-details-container Details">
  <button type="button" class="btn">Tab once from this button, and press enter</button>
  <button type="button" class="sr-only js-details-target">
    Screen reader only button
  </button>
  <div class="Details-content--hidden">
    Button activated!
  </div>
</div>
```

## The media object

Create a media object with utilities.

```html
<div class="clearfix p-3 border">
  <div class="float-left p-3 mr-3 bg-gray">
    Image
  </div>
  <div class="overflow-hidden">
    <p><b>Body</b> Bacon ipsum dolor amet shankle rump tenderloin pork chop meatball strip steak bresaola doner sirloin capicola biltong shank pig. Alcatra frankfurter ham hock, ribeye prosciutto hamburger kevin brisket chuck burgdoggen short loin.</p>
  </div>
</div>
```
Create a double-sided media object for a container with a flexible center.

```html
<div class="clearfix p-3 border border-gray">
  <div class="float-left p-3 mr-3 bg-gray">
    Image
  </div>
  <div class="float-right p-3 ml-3 bg-gray">
    Image
  </div>
  <div class="overflow-hidden">
    <p><b>Body</b> Bacon ipsum dolor amet shankle rump tenderloin pork chop meatball strip steak bresaola doner sirloin capicola biltong shank pig. Alcatra frankfurter ham hock, ribeye prosciutto hamburger kevin brisket chuck burgdoggen short loin.</p>
  </div>
</div>
```

You can also create a media object with [flexbox utilities](./flexbox#media-object) instead of floats which can be useful for changing the vertical alignment.
