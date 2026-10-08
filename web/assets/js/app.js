// Floor-plan tiles open the room they point at, the ticker can be paused, and sparkles
// trail the mouse, MySpace-style.
(function () {
	function openRoom(id) {
		const room = id && document.getElementById(id)
		if (room && room.tagName === 'DETAILS') room.open = true
	}

	document.addEventListener('click', function (event) {
		const link = event.target.closest('a[data-room]')
		if (link) openRoom(link.dataset.room)
	})
	openRoom(location.hash.slice(1))

	const ticker = document.querySelector('.ticker')
	const toggle = ticker && ticker.querySelector('.ticker__toggle')
	if (toggle) {
		toggle.addEventListener('click', function () {
			const paused = ticker.classList.toggle('is-paused')
			toggle.textContent = paused ? 'Play' : 'Pause'
		})
	}

	if (!matchMedia('(pointer: fine)').matches || matchMedia('(prefers-reduced-motion: reduce)').matches) return

	const colours = ['#FF2E93', '#FFB3D9', '#E0A526', '#FFFFFF']
	let last = 0
	let alive = 0
	addEventListener('pointermove', function (event) {
		if (event.timeStamp - last < 35 || alive > 24) return
		last = event.timeStamp

		const sparkle = document.createElement('span')
		sparkle.className = 'trail'
		sparkle.style.left = event.clientX + 'px'
		sparkle.style.top = event.clientY + 'px'
		sparkle.style.background = colours[Math.floor(Math.random() * colours.length)]
		sparkle.style.setProperty('--drift', Math.round(Math.random() * 30 - 15) + 'px')
		sparkle.addEventListener('animationend', function () {
			sparkle.remove()
			alive--
		})
		alive++
		document.body.appendChild(sparkle)
	}, { passive: true })
})()
