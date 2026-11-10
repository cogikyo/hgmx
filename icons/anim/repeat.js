import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const loop = [el.querySelector('path[d^="M3 11"]'), el.querySelector('path[d^="M21 13"]')];
  animate(loop, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });

  const top = el.querySelector('path[d^="m17 2"]');
  animate(top, { translateX: [-4, 2, 0], opacity: [0, 1] }, { duration: τmid });

  const bottom = el.querySelector('path[d^="m7 22"]');
  animate(bottom, { translateX: [4, -2, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
