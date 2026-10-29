import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M21 11.5"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rings = el.querySelectorAll('path[d$="v3"]');
  animate(rings, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmid }) });
  const plus = [el.querySelector('path[d="M16 18h6"]'), el.querySelector('path[d="M19 15v6"]')];
  animate(plus, { scale: [0, 1.25, 1], rotate: [0, 180] }, { duration: τmid, delay: τmin });
}
