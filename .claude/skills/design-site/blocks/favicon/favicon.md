# Logo and favicon

Every site gets a logo, drawn as one SVG from the attendee's brief. The same file is the header
logo and the browser-tab icon, so there's nothing to generate and no API key.

| File | Where |
|------|-------|
| the logo | `web/assets/images/logo.svg` |
| the manifest | `web/assets/manifest.json` |
| the tags below | the layout's `<head>` |

## Drawing the logo

Ask what it should be before drawing: an animal, a letter, a shape from their brief. Then write
the SVG by hand. It's a small file, so trying a few is cheap.

- A square `viewBox="0 0 64 64"` for the mark, since browsers show the favicon square.
- Shapes and paths only, no `<text>`: a favicon can't load web fonts, so text falls back to
  whatever the browser has. If they want a wordmark, set the name in the header as HTML text
  next to the mark, in the site's font.
- Two or three flat colours from their palette, and bold shapes. Check it still reads at 16px:
  thin lines and small details disappear.
- For a mark that suits dark tabs too, put a `<style>` inside the SVG with
  `@media (prefers-color-scheme: dark) { … }` that changes the fills.
- No scripts, no external references, no embedded raster images.

In the header, link it home with its own alt text:

```templ
<a class="logo" href="/">
	<img src="/assets/images/logo.svg" alt="Site name, home" width="40" height="40"/>
</a>
```

## Head tags

Straight after the viewport `<meta>`:

```templ
<link rel="icon" href="/assets/images/logo.svg" type="image/svg+xml"/>
<link rel="manifest" href="/assets/manifest.json"/>
<meta name="theme-color" content="#…"/>
```

`theme-color` is the page background, and it tints the browser's toolbar on phones.

## Manifest

It's named `manifest.json` rather than `.webmanifest` because Go's file server knows the `.json`
type.

```json
{
	"name": "Their site's name",
	"short_name": "Short name",
	"icons": [{ "src": "/assets/images/logo.svg", "sizes": "any", "type": "image/svg+xml" }],
	"theme_color": "#…",
	"background_color": "#…",
	"display": "standalone"
}
```

Keep `short_name` to 12 characters or fewer: it's the label under a home-screen icon.

## iPhone home screen (optional)

Safari's home-screen icon needs a PNG. If `rsvg-convert` or `magick` is installed, make one on the
page background, then add a tag for it:

```sh
rsvg-convert -w 180 -h 180 -b '#…' web/assets/images/logo.svg -o web/assets/images/apple-touch-icon.png
# or
magick -background '#…' -density 384 web/assets/images/logo.svg -resize 180x180 web/assets/images/apple-touch-icon.png
```

```templ
<link rel="apple-touch-icon" href="/assets/images/apple-touch-icon.png"/>
```

If neither tool is there, skip it and say so. Everything else works without it.

Adapted from needypeanut's `blocks/favicon`, which generates a full PNG set through the Real
Favicon Generator API. That needs a key, so this workshop version uses the SVG.
