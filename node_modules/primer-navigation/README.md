# Primer CSS Navigation

[![npm version](http://img.shields.io/npm/v/primer-navigation.svg)](https://www.npmjs.org/package/primer-navigation)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Primer comes with several navigation components. Some were designed with singular purposes, while others were design to be more flexible and appear quite frequently.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-navigation` with this command.

```
$ npm install --save primer-navigation
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-navigation/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Navigation
status: Stable
-->

Primer comes with several navigation components. Some were designed with singular purposes, while others were design to be more flexible and appear quite frequently.

{:toc}

## Menu

The menu is a vertical list of navigational links. **A menu's width and placement must be set by you.** If you like, just use our grid columns as a parent. Otherwise, apply a custom `width`.

```html
<nav class="menu" aria-label="Person settings">
  <a class="menu-item selected" href="#url" aria-current="page">Account</a>
  <a class="menu-item" href="#url">Profile</a>
  <a class="menu-item" href="#url">Emails</a>
  <a class="menu-item" href="#url">Notifications</a>
</nav>
```

There are a few subcomponents and add-ons that work well with the menu, including avatars, counters, and Octicons.

```html
<nav class="menu" aria-label="Person settings">
  <a class="menu-item selected" href="#url" aria-current="page">
    <%= octicon "tools" %>
    Account
  </a>
  <a class="menu-item" href="#url">
    <%= octicon "person" %>
    Profile
  </a>
  <a class="menu-item" href="#url">
    <%= octicon "mail" %>
    Emails
  </a>
  <a class="menu-item" href="#url">
    <%= octicon "radio-tower" %>
    <span class="Counter">3</span>
    Notifications
  </a>
</nav>
```

You can also add optional headings to a menu. Feel free to use nearly any semantic element with the `.menu-heading` class, including inline elements, headings, and more.

```html
<nav class="menu" aria-labelledby="menu-heading">
  <span class="menu-heading" id="menu-heading">Menu heading</span>
  <a class="menu-item selected" href="#url" aria-current="page">Account</a>
  <a class="menu-item" href="#url">Profile</a>
  <a class="menu-item" href="#url">Emails</a>
  <a class="menu-item" href="#url">Notifications</a>
</nav>
```


## Tabnav

When you need to toggle between different views, consider using a tabnav. It'll given you a left-aligned horizontal row of... tabs!

```html
<div class="tabnav">
  <nav class="tabnav-tabs" aria-label="Foo bar">
    <a href="#url" class="tabnav-tab selected" aria-current="page">Foo tab</a>
    <a href="#url" class="tabnav-tab">Bar tab</a>
  </nav>
</div>
```

Use `.float-right` to align additional elements in the `.tabnav`:

```html
<div class="tabnav">
  <a class="btn btn-sm float-right" href="#url" role="button">Button</a>
  <nav class="tabnav-tabs" aria-label="Foo bar">
    <a href="#url" class="tabnav-tab selected" aria-current="page">Foo Tab</a>
    <a href="#url" class="tabnav-tab">Bar Tab</a>
  </nav>
</div>
```

Additional bits of text and links can be styled for optimal placement with `.tabnav-extra`:

```html
<div class="tabnav">
  <div class="tabnav-extra float-right">
    Tabnav widget text here.
  </div>
  <nav class="tabnav-tabs" aria-label="Foo bar">
    <a href="#url" class="tabnav-tab selected" aria-current="page">Foo Tab</a>
    <a href="#url" class="tabnav-tab">Bar Tab</a>
  </nav>
</div>
```

```html
<div class="tabnav">
  <div class="float-right">
    <a class="tabnav-extra" href="#url">
      Tabnav extra link
    </a>
    <a class="tabnav-extra" href="#url">
      Tabnav extra link
    </a>
  </div>
  <nav class="tabnav-tabs" aria-label="Foo bar">
    <a href="#url" class="tabnav-tab selected" aria-current="page">Foo Tab</a>
    <a href="#url" class="tabnav-tab">Bar Tab</a>
  </nav>
</div>
```

## Filter list

A vertical list of filters. Grey text on white background. Selecting a filter from the list will fill its background with blue and make the text white.

```html
<ul class="filter-list">
  <li>
    <a href="#url" class="filter-item selected" aria-current="page">
      <span class="count" title="results">21</span>
      First filter
    </a>
  </li>
  <li>
    <a href="#url" class="filter-item">
      <span class="count" title="results">3</span>
      Second filter
    </a>
  </li>
  <li>
    <a href="#url" class="filter-item">
      Third filter
    </a>
  </li>
</ul>
```

## Sub navigation

`.subnav` is navigation that is typically used when on a dashboard type interface with another set of navigation above it. This helps distinguish navigation hierarchy.

```html
<nav class="subnav" aria-label="Respository">
  <a href="#url" class="subnav-item selected" aria-current="page">Item 1</a>
  <a href="#url" class="subnav-item">Item 2</a>
  <a href="#url" class="subnav-item">Item 3</a>
</nav>
```

You can have `subnav-search` in the subnav bar.

```html
<div class="subnav">
  <nav class="subnav-links" aria-label="Repository">
    <a href="#url" class="subnav-item selected" aria-current="page">Item 1</a>
    <a href="#url" class="subnav-item">Item 2</a>
    <a href="#url" class="subnav-item">Item 3</a>
  </nav>
  <div class="subnav-search float-left">
    <input type="search" name="name" class="form-control subnav-search-input" value="" aria-label="Search site">
    <%= octicon "search", :class => "subnav-search-icon" %>
  </div>
</div>
```


You can also use a `subnav-search-context` to display search help in a select menu.

```html
<div class="subnav">
  <nav class="subnav-links">
    <a href="#url" class="subnav-item selected">Item 1</a>
    <a href="#url" class="subnav-item">Item 2</a>
    <a href="#url" class="subnav-item">Item 3</a>
  </nav>
  <div class="float-left ml-3 select-menu js-menu-container js-select-menu subnav-search-context">
    <button type="button" name="button" class="btn select-menu-button js-menu-target" aria-expanded="false" aria-haspopup="true">Filters </button>
    <div class="select-menu-modal-holder js-menu-content js-navigation-container" aria-hidden="true">
      <div class="select-menu-modal">
        <div class="select-menu-list">
          <a href="#url" class="select-menu-item js-navigation-item">
            <div class="select-menu-item-text">
              Search filter 1
            </div>
          </a>
          <a href="#url" class="select-menu-item js-navigation-item">
            <div class="select-menu-item-text">
              Search filter 2
            </div>
          </a>
          <a href="#url" class="select-menu-item js-navigation-item">
            <div class="select-menu-item-text">
              Search filter 3
            </div>
          </a>
        </div>
      </div>
    </div>
  </div>
  <div class="subnav-search float-left">
    <input type="search" name="name" class="form-control subnav-search-input" value="" aria-label="Search site">
    <%= octicon "search", :class => "subnav-search-icon" %>
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
