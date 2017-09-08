# HEAD

# 9.3.0

## Added
- Docs for `primer-layout` (grid), `primer-support`, `primer-utilities`, and `primer-marketing-utilities`
- Primer keys for `category` and `module_type` to `package.json` (for use in documentation and gathering stats)

## Changes
- Removes `docs` from `gitignore`
- Removes the `^` from all dependencies so that we can publish exact versions
- Consolidates release notes from various sources into one changelog located in `/modules/primer-css/CHANGELOG.md`

# 9.2.0

## Added

- Add `test-docs` npm script in each module to check that every CSS class is documented (or at least mentioned) in the module's own markdown docs

## Changes

- Remove per-module configurations (`.gitignore`, `.postcss.json`, `.stylelintrc.json`) and `CHANGELOG.md` files in #284
- Replace most static `font-size`, `font-weight`, and `line-height` CSS property values with their [SCSS variable equivalents](https://github.com/primer/primer-css/blob/c9ea37316fbb73c4d9931c52b42bc197260c0bf6/modules/primer-support/lib/variables/typography.scss#L12-L33) in #252
- Refactor CI scripts to use Travis conditional deployment for release candidate and final release publish steps in #290

# 9.1.1

This release updates primer modules to use variables for spacing units instead of pixel values.

## Changes

- primer-alerts: 1.2.0 => 1.2.1
- primer-avatars: 1.1.0 => 1.1.1
- primer-base: 1.2.0 => 1.2.1
- primer-blankslate: 1.1.0 => 1.1.1
- primer-box: 2.2.0 => 2.2.1
- primer-breadcrumb: 1.1.0 => 1.1.1
- primer-buttons: 2.1.0 => 2.1.1
- primer-cards: 0.2.0 => 0.2.1
- primer-core: 6.1.0 => 6.1.1
- primer-css: 9.1.0 => 9.1.1
- primer-forms: 1.1.0 => 1.1.1
- primer-labels: 1.2.0 => 1.2.1
- primer-layout: 1.1.0 => 1.1.1
- primer-markdown: 3.4.0 => 3.4.1
- primer-marketing-type: 1.1.0 => 1.1.1
- primer-marketing-utilities: 1.1.0 => 1.1.1
- primer-marketing: 5.1.0 => 5.1.1
- primer-navigation: 1.1.0 => 1.1.1
- primer-page-headers: 1.1.0 => 1.1.1
- primer-page-sections: 1.1.0 => 1.1.1
- primer-product: 5.1.0 => 5.1.1
- primer-support: 4.1.0 => 4.1.1
- primer-table-object: 1.1.0 => 1.1.1
- primer-tables: 1.1.0 => 1.1.1
- primer-tooltips: 1.1.0 => 1.1.1
- primer-truncate: 1.1.0 => 1.1.1
- primer-utilities: 4.4.0 => 4.4.1

# 9.1.0

This release updates our [stylelint config](/primer/stylelint-config-primer) to [v2.0.0](https://github.com/primer/stylelint-config-primer/releases/tag/v2.0.0), and to stylelint v7.13.0. Each module also now has a `lint` npm script, and there are top-level `test` and `lint` scripts that you can use to lint and test all modules in one go.

This release also includes major improvements to our Travis build scripts to automatically publish PR builds, release candidates, and the "final" versions to npm.

# 9.0.0 - Core dependency & repo urls

We discovered that `primer-core` specified and outdated version of `primer-base` in it's dependencies. The outdated version did not have `normalize.scss` included which could cause some issues. This has issue occurred during v7.0.0 when creating the new monorepo. Also fixes repo urls in `package.json` for individual packages.

See PR [#243](https://github.com/primer/primer-css/pull/243)

## Changes

### Primer Core v6.0.0
- Fixed `primer-base` dependency to point to latest version

**Repo urls corrected from `packages` to `modules` in:**
- primer-base v1.1.5
- primer-box v2.1.8
- primer-buttons v2.0.6
- primer-forms v1.0.13
- primer-layout v1.0.5
- primer-navigation v1.0.6
- primer-support v4.0.7
- primer-table-object v1.0.9
- primer-tooltips v1.0.2
- primer-truncate v1.0.2
- primer-utilities v4.3.5

### Primer Product v5.0.2

**Repo urls corrected from `packages` to `modules` in:**
- primer-alerts v1.1.8
- primer-avatars v1.0.2
- primer-blankslate v1.0.2
- primer-labels v1.1.6
- primer-markdown v3.3.13
- primer-support v4.0.7

### Primer Marketing v5.0.2

**Repo urls corrected from `packages` to `modules` in:**
- primer-breadcrumb v1.0.2
- primer-cards v0.1.8
- primer-marketing-support v1.0.2
- primer-marketing-type v1.0.2
- primer-marketing-utilities v1.0.2
- primer-page-headers v1.0.2
- primer-page-sections v1.0.2
- primer-support v4.0.7
- primer-tables v1.0.2

# 8.0.0 - Imports

Fixes issues with the ordering of imports in each of our meta-packages. See PR [#239](https://github.com/primer/primer-css/pull/239)


## Changes

### Primer Core v5.0.1
- Re-ordered imports in `index.scss` to ensure utilities come last in the cascade

### Primer Product v5.0.1
- Re-ordered imports in `index.scss` to move markdown import to end of list to match former setup in GitHub.com

### Primer Marketing v5.0.1
- Re-ordered imports in `index.scss` to ensure marketing utilities come last in the cascade

# 7.0.0 - Monorepo

In an effort to improve our publishing workflow we turned Primer CSS into a monorepo, made this repo the source of truth for Primer by removing Primer modules from GitHub, and setup Lerna for managing multiple packages and maintaining independent versioning for all our modules.

This is exciting because:

- we can spend less time hunting down the cause of a broken build and more time focussing on making Primer more useful and robust for everyone to use
- we can be more confident that changes we publish won't cause unexpected problems on GitHub.com and many other GitHub websites that use Primer
- we no longer have files like package.json, scripts, and readme's in the GitHub app that don't really belong there
- **we can accept pull requests from external contributors** again!

See PR for more details on this change: https://github.com/primer/primer-css/pull/230

## Other changes:

### Primer Core v4.0.3

#### primer-support v4.0.5
- Update fade color variables to use rgba instead of transparentize color function for better Sass readability
- Update support variables and mixins to use new color variables

#### primer-layout v1.0.3
- Update grid gutter styles naming convention and add responsive modifiers
- Deprecate `single-column` and `table-column` from layout module
- Remove `@include clearfix` from responsive container classes

#### primer-utilities v4.3.3
- Add `show-on-focus` utility class for accessibility
- Update typography utilities to use new color variables
- Add `.p-responsive` class

#### primer-base v1.1.3
- Update `b` tag font weight to use variable in base styles

### Primer Marketing v4.0.3

#### primer-tables
- Update marketing table colors to use new variables


# 6.0.0
- Add `State--small` to labels module
- Fix responsive border utilities
- Added and updated typography variables and mixins; updated variables used in typography utilities; updated margin, padding, and typography readmes
- Darken `.box-shadow-extra-large` shadow
- Update `.tooltip-multiline` to remove `word-break: break-word` property
- Add `.border-purple` utility class
- Add responsive border utilities to primer-marketing
- Add `ws-normal` utility for `whitespace: normal`
- Updated syntax and classnames for `Counters` and `Labels`, moved into combined module with states.

# 5.1.0
- Add negative margin utilities
- Move `.d-flex` & `.d-flex-inline` to be with other display utility classes in `visibility-display.scss`
- Delete `.shade-gradient` in favor of `.bg-shade-gradient`
- Removed alt-body-font variable from primer-marketing
- Removed un-used `alt` typography styles from primer-marketing
- Add green border utility

# 5.0.0
- Added new border variable and utility, replaced deprecated flash border variables
- Updated variable name in form validation
- Updated `.sr-only` to not use negative margin
- Added and removed border variables and utilities
- Add filter utility to Primer Marketing
- Removed all custom color variables within Primer-marketing in favor of the new color system
- Updated style for form group error display so it is positioned properly
- Updated state closed color and text and background pending utilities
- Removed local font css file from primer-marketing/support
- Updated all color variables and replaced 579 hex refs across modules with new variables, added additional shades to start introducing a new color system which required updating nearly all primer modules
- Added layout utility `.sr-only` for creating screen reader only elements
- Added `.flex{-infix}-item-equal` utilities for creating equal width and equal height flex items.
- Added `.flex{-infix}-row-reverse` utility for reversing rows of content
- Updated `.select-menu-button-large` to use `em` units for sizing of the CSS triangle.
- Added `.box-shadow-extra-large` utility for large, diffused shadow
- Updated: removed background color from markdown body
- Updated: remove background on the only item in an avatar stack
- Added form utility `.form-checkbox-details` to allow content to be shown/hidden based on a radio button being checked
- Added form utility to override Webkit's incorrect assumption of where to try to autofill contact information

# 4.7.0
- Update primer modules to use bold variable applying `font-weight: 600`

# 4.6.0
- Added `CircleBadge` component for badge-like displays within product/components/avatars
- Added Box shadow utilities `box-shadow`, `box-shadow-medium`, `box-shadow-large`, `box-shadow-none`
- Moved visibility and display utilities to separate partial at the end of the imports list, moved flexbox to it's own partial
- Added `flex-shrink-0` to address Flexbox Safari bug
- Updated: Using spacing variables in the `.flash` component
- Updated Box component styles and documentation
- Added `.wb-break-all` utility

# 4.4.0
- Adding primer-marketing module to primer
- Added red and blue border color variables and utilities
- Updated: `$spacer-5` has been changed to `32px` from `36px`
- Updated: `$spacer-6` has been changed to `40px` from `48px`
- Deprecated `link-blue`, updated `link-gray` and `link-gray-dark`, added `link-hover-blue` - Updated: blankslate module to use support variables for sizing

# 4.3.0
- Renamed `.flex-table` to `.TableObject`
- Updated: `$spacer-1` has been changed to `4px` from `3px`
- Updated: `$spacer-2` has been changed to `6px` from `8px`
- Added: `.text-shadow-dark` & `.text-shadow-light` utilities
- Updated: Moved non-framework CSS out of Primer modules. Added `box.scss` to `primer-core`. Added `discussion-timeline.scss` to `primer-product`, and moved `blob-csv.scss` into `/primer-product/markdown` directory
- Added: Flex utilities
- Refactor: Site typography to use Primer Marketing styles
- Added: `.list-style-none` utility
- Refactor: Button groups into some cleaner CSS
- Updated: Reorganizing how we separate primer-core, primer-product, primer-marketing css


# 4.2.0
- Added: Responsive styles for margin and padding utilities, display,  float, and new responsive hide utility, and updates to make typography responsive
- Added: new container styles and grid styles with responsive options
- Added: updated underline nav styles
- Deprecate: Deprecating a lot of color and layout utilities
- Added: More type utilities for different weights and larger sizes.
- Added: Well defined browser support


# 4.1.0
- Added: [primer-markdown](https://github.com/primer/markdown) to the build
- Fixes: Pointing "style" package.json to `build/build.css` file.
- Added: Update font stack to system fonts
- Added: Updated type scale as part of system font update
- Added: `.Box` component for replacing boxed groups, simple box, and table-list styles
- Added: New type utilities for headings and line-height
- Deprecated: `vertical-middle` was replaced with `v-align-middle`.
- Added: Layout utilities for vertical alignment, overflow, width and height, visibility, and display table
- Added: Changing from font icons to SVG

# 4.0.2
- Added npm build scripts to add `build/build.css` to the npm package

# 4.0.1
- Fixed: missing primer-layout from build

# 4.0.0
- Whole new npm build system, pulling in the code from separate component repos

# 3.0.0
- Added: Animation utilities
- Added: Whitespace scale, and margin and padding utilities
- Added: Border utilities
