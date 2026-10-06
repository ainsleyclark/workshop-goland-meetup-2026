// block: swiper — copy to web/assets/js/swiper.js. See swiper.md.
//
// Starts every .carousel on the page. Each one reads its desktop column
// count from data-swiper-desktop, and the gutter from --swiper-gutter in
// swiper.css, so one script drives any number of them.

document.querySelectorAll('.carousel').forEach(function (el) {
	// Swiper comes from a CDN. Without it, swiper.css turns the row into a
	// plain scroll-snap strip.
	if (typeof Swiper === 'undefined') {
		el.classList.add('carousel--static');
		return;
	}

	// Unset means swiper.css's 24px fallback; an explicit 0px must stay 0.
	const raw = getComputedStyle(el).getPropertyValue('--swiper-gutter').trim();
	const gutter = raw === '' ? 24 : parseFloat(raw) || 0;
	const desktopSlides = parseInt(el.dataset.swiperDesktop, 10) || 3;
	el.style.setProperty('--desktop-cols', desktopSlides);

	new Swiper(el, {
		// The .15 peek is deliberate: it tells people the row swipes.
		slidesPerView: 1.15,
		spaceBetween: 16,
		slidesOffsetBefore: gutter,
		slidesOffsetAfter: gutter,
		pagination: { el: el.querySelector('.swiper-pagination'), clickable: true },
		breakpoints: {
			768: { slidesPerView: 2, spaceBetween: 20 },
			// Off at desktop: the grid in swiper.css takes over.
			1024: { enabled: false },
		},
	});
});
