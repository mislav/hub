---
title: Box shadow
status: New release
status_issue: https://github.com/github/design-systems/issues/269
---

Box shadows are used to make content appear elevated. They are typically applied to containers of content that users need to pay attention to or content that appears on top of (overlapping) other elements on the page (like a pop-over or modal).

{:toc}

## Default

Default shadows are mainly used on things that need to appear slightly elevated, like pricing cards or UI elements containing important information.

```html
<div class="box-shadow p-3">
  .box-shadow
</div>
```

These types of shadows are typically applied to elements with borders, like [`Box`](/styleguide/css/styles/core/components/box).

```html
<div class="col-5">
  <div class="Box box-shadow">
    <div class="Box-row">
      <h3 class="mb-0">Organization</h3>
    </div>
    <div class="Box-row">
      <p class="mb-0 alt-text-small text-gray">
        Taxidermy live-edge mixtape, keytar tumeric locavore meh selvage deep v letterpress vexillologist lo-fi tousled church-key thundercats. Brooklyn bicycle rights tousled, marfa actually.
      </p>
    </div>
    <div class="Box-row">
      <button type="button" name="Create an organization" class="btn btn-primary btn-block">Create an organization</button>
    </div>
  </div>
</div>
```

## Medium

Medium box shadows are more diffused and slightly larger than small box shadows.

```html
<div class="box-shadow-medium p-3">
  .box-shadow-medium
</div>
```

Medium box shadows are typically used on editorialized content that needs to appear elevated. Most of the time, the elements using this level of shadow will be clickable.

```html
<div class="col-6">
  <a class="d-block box-shadow-medium px-3 pt-4 pb-6 position-relative rounded-1 overflow-hidden no-underline" href="#">
    <div class="bg-blue position-absolute top-0 left-0 py-1 width-full"></div>
    <h3 class="text-gray-dark">Serverless architecture</h3>
    <p class="alt-text-small text-gray">
      Build powerful, event-driven, serverless architectures with these open-source libraries and frameworks.
    </p>
    <ul class="position-absolute bottom-0 pb-3 text-small text-gray list-style-none ">
      <li class="d-inline-block mr-1"><%= octicon "repo", :class => "mr-1" %>12 Repositories</li>
      <li class="d-inline-block"><%= octicon "code", :class => "mr-1" %>5 Languages</li>
    </ul>
  </a>
</div>
```

## Large

Large box shadows are very diffused and used sparingly in the product UI.

```html
<div class="box-shadow-large p-3">
  .box-shadow-large
</div>
```

These shadows are used for marketing content, UI screenshots, and content that appears on top of other page elements.

```html
<div class="box-shadow-large rounded-2 overflow-hidden">
  <img src="<%= image_path "modules/site/org_example_nasa.png" %>" class="img-responsive" alt="NASA is on GitHub">
</div>
```

## Extra Large

Extra large box shadows are even more diffused.

```html
<div class="box-shadow-extra-large p-3">
  .box-shadow-extra-large
</div>
```

These shadows are used for marketing content and content that appears on top of other page elements.

## Remove a box shadow

Additionally there is a `box-shadow-none` class that removes `box-shadow`:

```html
<div class="box-shadow-large box-shadow-none p-3">
  .box-shadow-none
</div>
```
