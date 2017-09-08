# Primer CSS Labels

[![npm version](http://img.shields.io/npm/v/primer-labels.svg)](https://www.npmjs.org/package/primer-labels)
[![Build Status](https://travis-ci.org/primer/primer-css.svg?branch=master)](https://travis-ci.org/primer/primer-css)

> Labels add metadata or indicate status of items and navigational elements.

This repository is a module of the full [primer-css][primer-css] repository.

## Install

This repository is distributed with [npm][npm]. After [installing npm][install-npm], you can install `primer-labels` with this command.

```
$ npm install --save primer-labels
```

## Usage

The source files included are written in [Sass][sass] (`scss`) You can simply point your sass `include-path` at your `node_modules` directory and import it like this.

```scss
@import "primer-labels/index.scss";
```

You can also import specific portions of the module by importing those partials from the `/lib/` folder. _Make sure you import any requirements along with the modules._

## Build

For a compiled **css** version of this module, a npm script is included that will output a css version to `build/build.css` The built css file is also included in the npm package.

```
$ npm run build
```

## Documentation

<!-- %docs
title: Labels
status_issue: https://github.com/github/design-systems/issues/332
status: New release
-->

Labels add metatdata or indicate status of items and navigational elements. Three different types of labels are available: [Labels](#default-label-styles) for adding metadata, [States](#states) for indicating status, and [Counters](#counters) for showing the count for a number of items.


{:toc}

## Labels

The base label component styles the text, adds padding and rounded corners, and an inset box shadow. Labels come in various themes which apply colors and different border styles.

GitHub also programmatically generates and applies a background color for labels on items such as issues and pull requests. Users are able to select any background color and the text color will adjust to work with light and dark background colors.

The base `Label` style does not apply a background color, here's an example using the `bg-blue` utility to apply a blue background:

```html
<span title="Label: default label" class="Label bg-blue">default label</span>
```

**Note:** Be sure to include a title attribute on labels, it's helpful for people using screen-readers to differentiate a label from other text. I.e. without the title attribute, the following example would read as _"New select component design"_, rather than identifying `design` as a label.

```html
<a href="#url">New select component</a><span class="Label bg-blue ml-1">design</span>
```

### Label themes

Labels come in a few different themes. Use a theme that helps communicate the content of the label, and ensure it's used consistently.

Use `Label--gray` to create a label with a light gray background and gray text. This label is neutral in color and can be used in contexts where all you need to communicate is metadata, or whe you want a label to feel less prominent compared with labels with stronger colors.

```html
<span title="Label: gray label" class="Label Label--gray">gray label</span>
```

Use `Label--gray-darker` to create a label with a dark-gray background color. This label is also neutral in color, however, since it's background is darker it can stand out more compared to `Label--gray`.

```html
<span title="Label: dark gray label" class="Label Label--gray-darker">dark gray label</span>
```

Use `Label--orange` to communicate "warning". The orange background color is very close to red, so avoid using next to labels with a red background color since most people will find it hard to tell the difference.

```html
<span title="Label: orange label" class="Label Label--orange">orange label</span>
```

Use `Label--outline` to create a label with gray text, a gray border, and a transparent background. The outline reduces the contrast of this label in combination with filled labels. Use this in contexts where you need it to stand out less than other labels and communicate a neutral message.

```html
<span title="Label: outline label" class="Label Label--outline">outlined label</span>
```

Use `Label--outline-green` in combination with `Label--outline` to communicate a positive message.

```html
<span title="Label: green outline label" class="Label Label--outline Label--outline-green">green outlined label</span>
```


## States

Use state labels to inform users of an items status. States are large labels with bolded text. The default state has a gray background.

```html
<span class="State">Default</span>
```

### State themes
States come in a few variations that apply different colors. Use the state that best communicates the status or function.

```html
<span title="Status: open" class="State State--green"><%= octicon "git-pull-request" %> Open</span>
<span title="Status: closed" class="State State--red"><%= octicon "git-pull-request" %> Closed</span>
<span title="Status: merged" class="State State--purple"><%= octicon "git-merge" %> Merged</span>
```

**Note:** Similar to [labels](#labels), you should include the title attribute on states to differentiate them from other content.

### Small states
Use `State--small` for a state label with reduced padding a smaller font size. This is useful in denser areas of content.

```html
<span title="Status: open" class="State State--green State--small"><%= octicon "issue-opened" %> Open</span>
<span title="Status: closed" class="State State--red State--small"><%= octicon "issue-closed" %> Closed</span>
```

## Counters

Use the `Counter` component to add a count to navigational elements and buttons. Counters come in 3 variations: the default `Counter` with a light gray background, `Counter--gray` with a dark-gray background and inverse white text, and `Counter--gray-light` with a light-gray background and dark gray text.

```html
<span class="Counter">16</span>
<span class="Counter Counter--gray">32</span>
<span class="Counter Counter--gray-light">64</span>
```

Use the `Counter` in navigation to indicate the number of items without the user having to click through or count the items, such as open issues in a GitHub repo. See more options in [navigation](../../core/components/navigation).

```html
<div class="tabnav">
  <nav class="tabnav-tabs" aria-label="Foo bar">
    <a href="#url" class="tabnav-tab selected" aria-current="page">Foo tab <span class="Counter">23</a>
    <a href="#url" class="tabnav-tab">Bar tab</a>
  </nav>
</div>
```

Counters can also be used in `Box` headers to indicate the number of items in a list. See more on the [box component](../../core/components/box).

```html
<div class="Box">
  <div class="Box-header">
    <h3 class="Box-title">
      Box title
      <span class="Counter Counter--gray">3</span>
    </h3>
  </div>
  <ul>
    <li class="Box-row">
      Box row one
    </li>
    <li class="Box-row">
      Box row two
    </li>
    <li class="Box-row">
      Box row three
    </li>
  </ul>
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
