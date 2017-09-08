# Primer CSS Buttons

[![npm version](http://img.shields.io/npm/v/primer-buttons.svg)](https://www.npmjs.org/package/primer-buttons)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Buttons are used for actions, like in forms, while textual hyperlinks are used for destinations, or moving from one page to another.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-buttons` with this command.

```
$ npm install --save primer-buttons
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-buttons/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Buttons
status: Stable
-->

Buttons are used for **actions**, like in forms, while textual hyperlinks are used for **destinations**, or moving from one page to another.

#### Buttons

```html
<button class="m-1 btn" type="button">Close</button>
<button class="m-1 btn btn-primary" type="button">Comment</button>
<button class="m-1 btn btn-danger" type="button">Rename repository</button>
<button class="m-1 btn btn-blue" type="button">Sign up for free</button>
<button class="m-1 btn btn-purple" type="button">Try if for free</button>
<br>
<button class="m-1 btn hover" type="button">Close</button>
<button class="m-1 btn btn-primary hover" type="button">Comment</button>
<button class="m-1 btn btn-danger hover" type="button">Rename repository</button>
<button class="m-1 btn btn-blue hover" type="button">Sign up for free</button>
<button class="m-1 btn btn-purple hover" type="button">Try if for free</button>
<br>
<button class="m-1 btn focus" type="button">Close</button>
<button class="m-1 btn btn-primary focus" type="button">Comment</button>
<button class="m-1 btn btn-danger focus" type="button">Rename repository</button>
<button class="m-1 btn btn-blue focus" type="button">Sign up for free</button>
<button class="m-1 btn btn-purple focus" type="button">Try if for free</button>
<br>
<button class="m-1 btn selected" type="button">Close</button>
<button class="m-1 btn btn-primary selected" type="button">Comment</button>
<button class="m-1 btn btn-danger selected" type="button">Rename repository</button>
<button class="m-1 btn btn-blue selected" type="button">Sign up for free</button>
<button class="m-1 btn btn-purple selected" type="button">Try if for free</button>
<br>
<button class="m-1 btn disabled" type="button">Close</button>
<button class="m-1 btn btn-primary disabled" type="button">Comment</button>
<button class="m-1 btn btn-danger disabled" type="button">Rename repository</button>
<button class="m-1 btn btn-blue disabled" type="button">Sign up for free</button>
<button class="m-1 btn btn-purple disabled" type="button">Try if for free</button>
```

#### Default buttons

Use the standard—yet classy—`.btn` for form actions and primary page actions. These are used extensively around the site.

When using a `<button>` element, **always specify a `type`**. When using a `<a>` element, **always add `role="button"` for accessibility**.

```html
<button class="btn" type="button">Button button</button>
<a class="btn" href="#url" role="button">Link button</a>
```

You can find them in two sizes: the default `.btn` and the smaller `.btn-sm`.

```html
<button class="btn" type="button">Button</button>
<button class="btn btn-sm" type="button">Small button</button>
```

#### Primary

Primary buttons are green and are used to indicate the *primary* action on a page. When you need your buttons to stand out, use `.btn.btn-primary`. You can use it with both button sizes—just add `.btn-primary`.

```html
<button class="btn btn-primary" type="button">Primary button</button>
<button class="btn btn-sm btn-primary" type="button">Small primary button</button>
```

#### Danger

Danger buttons are red. They help reiterate that the intended action is important or potentially dangerous (e.g., deleting a repo or transferring ownership). Similar to the primary buttons, just add `.btn-danger`.

```html
<button class="btn btn-danger" type="button">Danger button</button>
<button class="btn btn-sm btn-danger" type="button">Small danger button</button>
```

#### Outline

Outline buttons downplay an action as they appear like boxy links. Just add `.btn-outline` and go.

```html
<button class="btn btn-outline" type="button">Outline button</button>
<button class="btn btn-sm btn-outline" type="button">Outline button</button>
```

#### Disabled state

Disable `<button>` elements with the boolean `disabled` attribute and `<a>` elements with the `.disabled` class.

```html
<button class="btn" type="button" disabled>Disabled button</button>
<a class="btn disabled" href="#url" role="button">Disabled button</a>
```

Similar styles are applied to primary, danger, and outline buttons:

```html
<button class="btn btn-primary" type="button" disabled>Disabled button</button>
<a class="btn btn-primary disabled" href="#url" role="button">Disabled button</a>
```

```html
<button class="btn btn-danger" type="button" disabled>Disabled button</button>
<a class="btn btn-danger disabled" href="#url" role="button">Disabled button</a>
```

```html
<button class="btn btn-outline" type="button" disabled>Disabled button</button>
<a class="btn btn-outline disabled" href="#url" role="button">Disabled button</a>
```

#### Block buttons

Make any button full-width by adding `.btn-block`. It adds `width: 100%;`, changes the `display` from `inline-block` to `block`, and centers the button text.

```html
<p><button class="btn btn-block" type="button">Block button</button></p>
<p><button class="btn btn-sm btn-block" type="button">Small block button</button></p>
```

#### With counts

You can easily append a count to a **small button**. Add the `.with-count` class to the `.btn-sm` and then add the `.social-count` after the button.

**Be sure to clear the float added by the additional class.**

```html
<div class="clearfix">
  <a class="btn btn-sm btn-with-count" href="#url" role="button">
    <%= octicon "eye" %>
    Watch
  </a>
  <a class="social-count" href="#url">6</a>
</div>
```

You can also use the [counter](../../product/components/labels) component within buttons:

```html
<button class="btn" type="button">
  Button
  <span class="Counter">12</span>
</button>

<button class="btn btn-primary" type="button">
  Button
  <span class="Counter">12</span>
</button>

<button class="btn btn-danger" type="button">
  Button
  <span class="Counter">12</span>
</button>

<button class="btn btn-outline" type="button">
  Button
  <span class="Counter">12</span>
</button>
```

#### Button groups

Have a hankering for a series of buttons that are attached to one another? Wrap them in a `.BtnGroup` and the buttons will be rounded and spaced automatically.

```html
<div class="BtnGroup mr-2">
  <button class="btn BtnGroup-item" type="button">Button</button>
  <button class="btn BtnGroup-item" type="button">Button</button>
  <button class="btn BtnGroup-item" type="button">Button</button>
</div>

<div class="BtnGroup mr-2">
  <button class="btn BtnGroup-item btn-outline" type="button">Button</button>
  <button class="btn BtnGroup-item btn-outline" type="button">Button</button>
  <button class="btn BtnGroup-item btn-outline" type="button">Button</button>
</div>

<div class="BtnGroup">
  <button class="btn BtnGroup-item btn-sm" type="button">Button</button>
  <button class="btn BtnGroup-item btn-sm" type="button">Button</button>
  <button class="btn BtnGroup-item btn-sm" type="button">Button</button>
</div>
```

Add `.BtnGroup-form` to `<form>`s within `.BtnGroup`s for proper spacing and rounded corners.

```html
<div class="BtnGroup">
  <button class="btn BtnGroup-item" type="button">Button</button>
  <form class="BtnGroup-form">
    <button class="btn BtnGroup-item" type="button">Button in a form</button>
  </form>
  <button class="btn BtnGroup-item" type="button">Button</button>
  <button class="btn BtnGroup-item" type="button">Button</button>
</div>
```

#### Hidden text button

Use `.hidden-text-expander` to indicate and toggle hidden text.

```html
<span class="hidden-text-expander">
  <button type="button" class="ellipsis-expander" aria-expanded="false">&hellip;</button>
</span>
```

You can also make the expander appear inline by adding `.inline`.

<!-- %enddocs -->

## License

[MIT](./LICENSE) &copy; [GitHub](https://github.com/)

[primer-css]: https://github.com/primer/primer
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
