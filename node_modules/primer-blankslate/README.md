# Primer CSS Blankslate

[![npm version](http://img.shields.io/npm/v/primer-blankslate.svg)](https://www.npmjs.org/package/primer-blankslate)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Blankslates are for when there is a lack of content within a page or section. Use them as placeholders to tell users why something isn’t there. Be sure to provide an action to add content as well.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-blankslate` with this command.

```
$ npm install --save primer-blankslate
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-blankslate/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Blankslate
status: Stable
-->

Blankslates are for when there is a lack of content within a page or section. Use them as placeholders to tell users why something isn't there. Be sure to provide an action to add content as well.

#### Basic example

Wrap some content in the outer `.blankslate` wrapper to give it the blankslate appearance.

```html
<div class="blankslate">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

#### With Octicons

When it helps the message, include (relevant) icons in your blank slate. Add `.blankslate-icon` to any `.mega-octicon`s as the first elements in the blankslate, like so.

```html
<div class="blankslate">
  <%= octicon "git-commit", :height => 32, :class => "blankslate-icon" %>
  <%= octicon "tag", :height => 32, :class => "blankslate-icon" %>
  <%= octicon "git-branch", :height => 32, :class => "blankslate-icon" %>
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

#### Variations

Add an additional optional class to the `.blankslate` to change its appearance.

##### Narrow

Narrows the blankslate container to not occupy the entire available width.

```html
<div class="blankslate blankslate-narrow">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

##### Capped

Removes the `border-radius` on the top corners.

```html
<div class="blankslate blankslate-capped">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

##### Spacious

Significantly increases the vertical padding.

```html
<div class="blankslate blankslate-spacious">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

##### Large

Increases the size of the text in the blankslate

```html
<div class="blankslate blankslate-large">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
</div>
```

##### No background

Removes the `background-color` and `border`.

```html
<div class="blankslate blankslate-clean-background">
  <h3>This is a blank slate</h3>
  <p>Use it to provide information when no dynamic content exists.</p>
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
