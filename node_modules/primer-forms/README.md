# Primer CSS Forms

[![npm version](http://img.shields.io/npm/v/primer-forms.svg)](https://www.npmjs.org/package/primer-forms)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Style individual form controls and utilize common layouts.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-forms` with this command.

```
$ npm install --save primer-forms
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-forms/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Forms
status: Stable
-->

Style individual form controls and utilize common layouts.

- We reset several elements' default styles for cross browser consistency and easier manipulation later. This includes `<fieldset>`s, WebKit validation bubbles, and most textual `<input>`s.
- Specific types of textual `<input>`s are styled automatically, but `.form-control` is available should you need it.
- Always declare a `type` on your `<button>`s.
- Form layouts rely on form groups.

#### Kitchen sink

All our inputs and buttons side-by-side for easy testing of sizing and alignment with one another.

```html
<p>
  <button class="btn" type="button">Button</button>
  <button class="btn select-menu-button" type="button" aria-expanded="false" aria-haspopup="true">
    Select menu
  </button>
  <input class="form-control" type="text" placeholder="Standard input" aria-label="Repository description">
  <input class="form-control input-monospace" type="text" placeholder="Monospace input" aria-label="SHA">
  <select class="form-select" aria-label="Important decision">
    <option>Select</option>
    <option value="option 2">Option 2</option>
  </select>
</p>

<p>
  <button class="btn btn-sm" type="button">Small button</button>
  <button class="btn btn-sm select-menu-button" type="button" aria-expanded="false" aria-haspopup="true">
    Select menu
  </button>
  <input class="form-control input-sm" type="text" placeholder="Small input" aria-label="Repository description">
  <select class="form-select select-sm" aria-label="Important decision">
    <option>Select</option>
    <option value="option 1">Option 1</option>
  </select>
</p>
```

#### Example form

Form controls in Primer currently have no basic layout specified (this is by design). You'll need to use `<fieldset>`s, `<div>`s, or other elements and styles to rearrange them.

```html
<form>
  <label for="name">Name</label>
  <input class="form-control" type="text" id="name">

  <label for="email">Email address</label>
  <input class="form-control" type="email" id="email">

  <label>
    <input type="checkbox"> Remember me
  </label>

  <label>
    <input type="radio" id="herp" name="herpderp" checked> Herp
  </label>
  <label>
    <input type="radio" id="derp" name="herpderp"> Derp
  </label>

  <button class="btn" type="submit">Submit</button>
</form>
```

#### Form contrast

Textual form controls have a white background by default. You can change this to a light gray with `.input-contrast`.

```html
<form>
  <input class="form-control" type="text" placeholder="Default input" aria-label="Default input">
  <input class="form-control input-contrast" type="text" placeholder="Input with contrast" aria-label="Input with contrast">
</form>
```

#### Sizing

Make inputs smaller, larger, or full-width with an additional class.

##### Small

```html
<form>
  <input class="form-control input-sm" type="text" placeholder="Small input" aria-label="Small input">
</form>
```

##### Large

```html
<form>
  <input class="form-control input-lg" type="text" placeholder="Large input" aria-label="Large input">
</form>
```

##### Block

```html
<form>
  <input class="form-control input-block" type="text" placeholder="Full-width input" aria-label="Full-width input">
</form>
```

##### Hide webkit's contact info autofill icon
Webkit sometimes gets confused and tries to add an icon/dropdown to autofill contact information on fields that may not be appropriate (such as input for number of users). Use this class to override the display of this icon.

```html
<form>
  <input class="form-control input-hide-webkit-autofill" placeholder="Hide Webkit's contact autofill on this input" type="text" aria-label="Hide Webkit's contact autofill on this input">
</form>
```

#### Selects

Primer adds light `height` and `vertical-align` styles to `<select>`s for all browsers to render them consistently with textual inputs.

```html
<form>
  <select class="form-select" aria-label="Preference">
    <option>Choose an option</option>
    <option>Git</option>
    <option>Subversion</option>
    <option>Social Coding</option>
    <option>Beets</option>
    <option>Bears</option>
    <option>Battlestar Galactica</option>
  </select>
</form>
```

##### Small

Use the `.select-sm` class to resize both default and custom `<select>`s to match the size of [our small buttons](/buttons/#default-buttons).

```html
<select class="form-select select-sm" aria-label="Preference">
  <option>Choose an option</option>
  <option>Git</option>
  <option>Subversion</option>
  <option>Social Coding</option>
  <option>Beets</option>
  <option>Bears</option>
  <option>Battlestar Galactica</option>
</select>

<select class="form-select select-sm" aria-label="Preference">
  <option>Choose an option</option>
  <option>Git</option>
  <option>Subversion</option>
  <option>Social Coding</option>
  <option>Beets</option>
  <option>Bears</option>
  <option>Battlestar Galactica</option>
</select>
```

#### Form groups

```html
<form>
  <dl class="form-group">
    <dt><label for="example-text">Example Text</label></dt>
    <dd><input class="form-control" type="text" value="Example Value" id="example-text"></dd>
  </dl>

  <dl class="form-group">
    <dt><label for="example-select">Example Select</label></dt>
    <dd>
      <select class="form-select" id="example-select">
        <option>Choose an option</option>
        <option>Git</option>
        <option>Subversion</option>
        <option>Social Coding</option>
        <option>Beets</option>
        <option>Bears</option>
        <option>Battlestar Galactica</option>
      </select>
    </dd>
  </dl>

  <dl class="form-group">
    <dt><label for="example-textarea">Example Textarea</label></dt>
    <dd>
      <textarea class="form-control" id="example-textarea"></textarea>
    </dd>
  </dl>
</form>
```

#### Form group validation

Convey errors and warnings for form groups. Add the appropriate class—either `.errored` or `.warn`—to the `<dl class="form-group">` to start. Then, house your error messaging in an additional `<dd>` with either `.error` or `.warning`.

```html
<form>
  <dl class="form-group errored">
    <dt><label for="example-text-errored">Example Text</label></dt>
    <dd><input class="form-control" type="text" value="Example Value" id="example-text-errored" aria-describedby="form-error-text"></dd>
    <dd class="error" id="form-error-text">This example input has an error.</dd>
  </dl>
  <br>
  <dl class="form-group warn">
    <dt><label for="example-text-warn">Example Text</label></dt>
    <dd><input class="form-control" type="text" value="Example Value" id="example-text-warn" aria-describedby="form-warning-text"></dd>
    <dd class="warning" id="form-warning-text">This example input has a warning.</dd>
  </dl>
</form>
```

#### Checkboxes and radios

Utilities to spice up the alignment and styling of checkbox and radio controls.

```html
<form>
  <div class="form-checkbox">
    <label>
      <input type="checkbox" checked="checked" aria-describedby="help-text-for-checkbox">
      Available for hire
    </label>
    <p class="note" id="help-text-for-checkbox">
      This will do insanely <strong>awesome</strong> and <strong>amazing</strong> things!
    </p>
  </div>
</form>
```

You may also add emphasis to the label:

```html
<form>
  <div class="form-checkbox">
    <label>
      <input type="checkbox" checked="checked">
      <em class="highlight">Available for hire</em>
    </label>
  </div>
</form>
```

##### Show / hide content based on a checkbox or radio button state

Content that is hidden by default should only be done so if it is non-essential for the context of the surrounding content. Be sure to use the `aria-live="polite"` attribute on the parent label for added content to be announced when displayed.

```html
<form>
  <div class="form-checkbox">
    <label>
      <input type="radio" name="hireme">
      Not available for hire
    </label>
  </div>
  <div class="form-checkbox">
    <label aria-live="polite">
      <input type="radio" class="form-checkbox-details-trigger" name="hireme">
      Available for hire
      <span class="form-checkbox-details text-normal">
        <span class="d-block mb-1">Available hours per week</span>
        <input type="text" name="" class="form-control input-contrast" size="3">
        <span class="text-small text-gray pl-2">hours per week</span>
      </span>
    </label>
  </div>
</form>
```

#### Input group

Attached an input and button to one another.

```html
<form>
  <div class="input-group">
    <input class="form-control" type="text" placeholder="Username" aria-label="Username">
    <span class="input-group-button">
      <button class="btn" type="button" aria-label="Copy to clipboard">
        <%= octicon "clippy" %>
      </button>
    </span>
  </div>
</form>
```

#### Form actions

Align buttons to the right—via `float: right` on the buttons—in forms with `.form-actions`. The floats are automatically cleared for you.

```html
<div class="form-actions">
  <button type="submit" class="btn btn-primary">Save changes</button>
  <button type="button" class="btn">Cancel</button>
</div>
```
<!-- %enddocs -->

## License

[MIT](./LICENSE) &copy; [GitHub](https://github.com/)

[primer-css]: https://github.com/primer/primer
[docs]: http://primercss.io/
[npm]: https://www.npmjs.com/
[install-npm]: https://docs.npmjs.com/getting-started/installing-node
[sass]: http://sass-lang.com/
