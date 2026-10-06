# Page transitions

A cross-fade when the reader moves between pages, with the header and logo staying where they
are. A few lines of CSS and no JavaScript: it's the
[View Transitions API](https://developer.mozilla.org/en-US/docs/Web/API/View_Transition_API) for
ordinary links between pages of the same site. Browsers without it navigate as they always did.

Worth adding once the site has a second page: on one page there's nothing to move between.

| File | Copy to |
|------|---------|
| `transitions.css` | `web/assets/css/transitions.css` |

## Wiring it up

In the layout's `<head>`, after the site's stylesheet:

```templ
<link rel="stylesheet" href="/assets/css/transitions.css"/>
```

Then match the two class names in the CSS to the layout: `.logo` is the favicon block's link, and
`.site-header` whatever wraps the header. Rename them if the layout uses other names.

## Gotchas

1. **Same site, whole pages.** It runs on `<a href>`s between pages this app serves, all on one
   origin, on `make web`'s localhost as much as the deployed site. It doesn't run on the first
   load, a reload, or a link away.
2. **One element per name per page.** If two elements share a `view-transition-name` the browser
   skips the transition altogether. So never name every card: thousands of names means thousands
   of snapshots per navigation. The morph from a card to its own page is a how-to below, with that
   in mind.
3. **Firefox and older Safari show nothing.** Chrome, Edge and Safari from 18.2 do it; the rest get
   a plain navigation, which is the fallback.
4. **Reduced motion is honoured.** The swap is instant.
5. **A slow page looks frozen.** The old page stays until the new one has arrived, so a page that
   takes a second to render spends that second doing nothing visible. Keep pages quick, or leave
   one out as below.

## Common edits

| Request | Change |
|---------|--------|
| Slide instead of fade | `::view-transition-old(root) { animation: 180ms ease-in both to-left }` and `::view-transition-new(root) { animation: 180ms ease-out both from-right }`, with two `@keyframes` moving `transform: translateX` and `opacity` |
| Keep the nav still too | `.site-nav { view-transition-name: site-nav }` |
| Morph a card into its page | the same name on the card and on the page's `<h1>`: `style={ map[string]string{"view-transition-name": "sighting-" + s.ID.String()} }`; the id is ours, not GBIF's text. Only on a page that shows a handful of cards, never the home page |
| No transition on one page | `<style>@view-transition { navigation: none; }</style>` in that page's `<head>` |
