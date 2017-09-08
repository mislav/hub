# Primer CSS Layout

[![npm version](http://img.shields.io/npm/v/primer-layout.svg)](https://www.npmjs.org/package/primer-layout)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Primer’s layout includes basic page containers and a single-tiered, fraction-based grid system. That sounds more complicated than it really is though—it’s just containers, rows, and columns.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-layout` with this command.

```
$ npm install --save primer-layout
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-layout/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Layout
status: Deprecated
status_issue: https://github.com/github/design-systems/issues/59
key: /css/styles/core/objects/layout
-->

Primer's layout includes basic page containers and a single-tiered, fraction-based grid system. That sounds more complicated than it really is though—it's just containers, rows, and columns.

You can find all the below styles in `_layout.scss`.

#### Container

Center your page's contents with a `.container`.

```html+erb
<div class="container">
  <!-- contents here -->
</div>
```

The container applies `width: 980px;` and uses horizontal `margin`s to center it.

#### Grid

##### How it works

The grid is pretty standard—you create rows with `.columns` and individual columns with a column class and fraction class. Here's how it works:

- Add a `.container` to encapsulate everything and provide ample horizontal gutter space.
- Create your outer row to clear the floated columns with `<div class="columns">`.
- Add your columns with individual `<div class="column">`s.
- Add your fractional width classes to set the width of the columns (e.g., `.one-fourth`).

##### Demo

In practice, your columns will look like the example below.

```html
<div class="container">
  <div class="columns mb-1">
    <div class="one-fifth column block-blue p-3 border">
      .one-fifth
    </div>
    <div class="four-fifths column block-blue p-3 border">
      .four-fifths
    </div>
  </div>

  <div class="columns mb-1">
    <div class="one-fourth column block-blue p-3 border">
      .one-fourth
    </div>
    <div class="three-fourths column block-blue p-3 border">
      .three-fourths
    </div>
  </div>

  <div class="columns mb-1">
    <div class="one-third column block-blue p-3 border">
      .one-third
    </div>
    <div class="two-thirds column block-blue p-3 border">
      .two-thirds
    </div>
  </div>

  <div class="columns">
    <div class="one-half column block-blue p-3 border">
      .one-half
    </div>
    <div class="one-half column block-blue p-3 border">
      .one-half
    </div>
  </div>
</div>
```

##### Centered

Columns can be [centered](/utilities/#centering-content) by adding `.centered` to the `.column` class.

```html
<div class="columns">
  <div class="one-half column centered block-blue p-3">
    .one-half
  </div>
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
