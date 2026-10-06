// block: threejs — copy to web/assets/js/scenes/<name>.js, then rewrite it. See threejs.md.
//
// The smallest scene that shows the contract: build in setup, move in the
// function it returns. Replace all of it with the attendee's own idea.

import { mount } from '/assets/js/threejs.js';

mount('example', function ({ THREE, scene, camera, colour, pointer }) {
	const shape = new THREE.Mesh(
		new THREE.IcosahedronGeometry(1.4, 0),
		new THREE.MeshStandardMaterial({ color: colour('--three-primary', '#2563eb'), flatShading: true }),
	);
	scene.add(shape);

	scene.add(new THREE.AmbientLight(0xffffff, 0.8));
	const key = new THREE.DirectionalLight(0xffffff, 2.5);
	key.position.set(3, 4, 5);
	scene.add(key);
	// A rim of the accent from behind, so the edges pick up the palette too.
	const rim = new THREE.DirectionalLight(colour('--three-accent', '#ffffff'), 3);
	rim.position.set(-4, 2, -3);
	scene.add(rim);

	camera.position.z = 5;

	return function (time) {
		shape.rotation.x = time * 0.2 + pointer.y * 0.4;
		shape.rotation.y = time * 0.3 + pointer.x * 0.6;
	};
});
