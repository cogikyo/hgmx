import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelectorAll('path[d^="M18 2"], path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = [el.querySelector("g"), el.querySelector('path[d^="m9 5"]')];
  animate(arrow, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
