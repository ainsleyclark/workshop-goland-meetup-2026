# Lightbox

Click a photo to see it full size over the page. Esc, the close button or a click outside the photo
closes it. It uses the browser's own `<dialog>`, so focus, Esc and the backdrop come for free, and
it needs no library.

| File | Copy to |
|------|---------|
| `lightbox.css` | `web/assets/css/lightbox.css` |
| `lightbox.js` | `web/assets/js/lightbox.js` |
| the component below | `web/views/components/lightbox.templ` |

## Component

One per page. It goes in the layout, after `{ children... }`:

```templ
package components

templ Lightbox() {
	<dialog id="lightbox" class="lightbox" aria-label="Photo">
		<form method="dialog">
			<button class="lightbox__close" aria-label="Close">×</button>
		</form>
		<figure>
			<img class="lightbox__image" alt=""/>
			<figcaption class="lightbox__caption" hidden></figcaption>
		</figure>
	</dialog>
}
```

`<form method="dialog">` closes the dialog without any script.

## Using it

Wrap each thumbnail in a link to the full photo, marked `data-lightbox`. Without JS, the link
simply opens the photo:

```templ
<a href={ templ.URL(m.URL) } data-lightbox data-caption={ m.Creator }>
	<img src={ m.URL } alt={ altText(m) } loading="lazy"/>
</a>
```

Load the two files in the layout:

```templ
<link rel="stylesheet" href="/assets/css/lightbox.css"/>
<script src="/assets/js/lightbox.js" defer></script>
```

## Gotchas

- **The photo URL and caption are GBIF's, so untrusted.** `templ.URL` sanitises the `href`, and
  the script only opens `http(s)` links and sets the caption with `textContent`. Keep it that way:
  never `innerHTML`.
- The lightbox reuses the thumbnail's `alt` for the big photo, so good alt text on thumbnails
  counts twice.
- It works on photos inside a swiper too: Swiper swallows the click at the end of a drag, so
  swiping doesn't open it.
- Style it from the site's palette: backdrop colour, caption font, the close button.
