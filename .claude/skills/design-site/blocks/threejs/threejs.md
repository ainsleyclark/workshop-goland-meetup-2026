# three.js

The plumbing for a 3D scene with [three.js](https://threejs.org). It draws nothing of its own:
it loads three.js without a bundler, sizes the canvas, pauses off screen, handles reduced motion
and falls back when WebGL isn't there. It uses three.js r186 from jsDelivr.

This shows what three.js could look like in the site, not how it must. What the scene is, where
it sits, how it looks and how the page around it is built are the attendee's to decide. Don't
offer ideas or steer them towards a kind of scene or layout: ask what they picture, then build
that, changing anything here that gets in the way.

| File | Copy to |
|------|---------|
| `threejs.js` | `web/assets/js/threejs.js` |
| `example.js` | `web/assets/js/scenes/<name>.js`, renamed and rewritten |
| `threejs.css` | `web/assets/css/threejs.css` |
| the component below | `web/views/components/three_scene.templ` |

## Wiring it up

In the layout's `<head>`, before any other module script. The import map is how three.js is
meant to load without a bundler, and it lets the addons find `three` too:

```templ
<script type="importmap">
	{
		"imports": {
			"three": "https://cdn.jsdelivr.net/npm/three@0.186.1/build/three.module.js",
			"three/addons/": "https://cdn.jsdelivr.net/npm/three@0.186.1/examples/jsm/"
		}
	}
</script>
<link rel="stylesheet" href="/assets/css/threejs.css"/>
```

## Component

```templ
package components

type ThreeSceneProps struct {
	Scene    string // A module in web/assets/js/scenes, without ".js".
	Label    string // What it shows, for screen readers; empty when it's decoration.
	Backdrop bool   // Fill the parent, behind its content.
}

templ ThreeScene(props ThreeSceneProps) {
	<canvas
		class={ "three-scene", templ.KV("three-scene--backdrop", props.Backdrop) }
		data-scene={ props.Scene }
		if props.Label != "" {
			role="img"
			aria-label={ props.Label }
		} else {
			aria-hidden="true"
		}
	></canvas>
	<script type="module" src={ "/assets/js/scenes/" + props.Scene + ".js" }></script>
}
```

Put it wherever their design wants it. The markup around it is theirs.

## Writing a scene

A scene is one module in `web/assets/js/scenes/`. It calls `mount` with the name it was given in
`ThreeSceneProps.Scene`. `setup` builds the scene once and returns a function that runs every
frame:

```js
import { mount } from '/assets/js/threejs.js';

mount('<name>', function ({ THREE, scene, camera, renderer, canvas, pointer, colour, still }) {
	// Build: meshes, lights, the camera's position.
	return function (time, delta) {
		// Move: time and delta are in seconds.
	};
});
```

| In `setup` | What it is |
|------------|------------|
| `THREE` | the three.js namespace, the same copy everywhere |
| `scene`, `camera`, `renderer` | ready to use: a `PerspectiveCamera(50)` at `z = 5`, and a transparent renderer sized to the parent |
| `canvas` | the `<canvas>` |
| `pointer` | `{ x, y }` from -1 to 1 across the window, for parallax or hover |
| `colour(prop, fallback)` | a CSS custom property off the canvas as a `THREE.Color`, so the scene uses the site's palette |
| `still` | true under reduced motion, when only one frame at `time = 0` is drawn |

Addons import from `three/addons/`, such as
`import { OrbitControls } from 'three/addons/controls/OrbitControls.js'`.

`example.js` is the smallest scene there is: a shape in `--three-primary`, rimmed with
`--three-accent`, that turns and leans towards the pointer. It only shows the contract, so its
shape isn't a suggestion: replace it with what they asked for rather than restyling it.

## Gotchas

1. **Restart `make web` after adding the component.** In development the JS is read from disk,
   but the templ is compiled in, so the page won't have the canvas until the server restarts.
2. **Load `three` only through the import map**, never from a second URL. Two copies of three.js
   on a page break each other.
3. **One or two scenes a page, never one per sighting.** The home page renders thousands of them,
   and each canvas is a WebGL context. To plot sightings, draw them all with one `Points` or
   `InstancedMesh`, with the data written out by `templ.JSONScript`. Send numbers and known
   fields only, and never put GBIF text into `innerHTML`.
4. **The parent's CSS background is the fallback.** The canvas stays transparent until it has
   drawn, so a blocked CDN or no WebGL just shows the background. With `Backdrop`, the canvas
   fills its parent, so the parent needs `position: relative` and the content in it
   `position: relative; z-index: 1`, or the canvas covers it. Without `Backdrop`, the canvas
   fills whatever box it's in, so that box needs a size.
5. **Decoration is `aria-hidden`.** If the scene means something, give it a `Label` and keep the
   same information in the HTML.
6. **Colours must be hex, `rgb()` or `hsl()`** for `colour()` to read them. `oklch()` and
   `color-mix()` don't parse.
7. **Models (`.glb`) go in `web/assets/models/`.** They're embedded in the binary, so keep them
   small, and load them with `GLTFLoader` from `three/addons/loaders/GLTFLoader.js`.
8. **Respect `still`.** Anything that moves without the frame function, such as controls that
   auto-rotate or a tween, should check it first.

## How-tos

Answers for when their idea needs them, not ideas to offer.

| To | Do |
|----|----|
| Let them drag to turn it | `OrbitControls` from `three/addons/controls/OrbitControls.js`, with `new OrbitControls(camera, canvas)`; without `Backdrop`, since that sets `pointer-events: none` |
| Show a model they supply | `GLTFLoader` from `three/addons/loaders/GLTFLoader.js`, then `loader.load('/assets/models/x.glb', (gltf) => scene.add(gltf.scene))` |
| Draw thousands of points | one `THREE.Points` with a `BufferGeometry` of positions, not a mesh per point |
| Get sharper edges | already on (`antialias: true`); the pixel ratio is capped at 1.5 in `threejs.js` |
| Keep it moving under reduced motion | don't; draw a pleasing first frame at `time = 0` instead |
| Plot sightings on a globe | a `SphereGeometry` for the earth, then one `Points` from the coordinates, written out with `templ.JSONScript` as `[lat, lng]` pairs: with both in radians, `x = r·cos(lat)·sin(lng)`, `y = r·sin(lat)`, `z = r·cos(lat)·cos(lng)` |
