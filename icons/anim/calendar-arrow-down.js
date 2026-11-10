import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M21 10.354"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const arrow = el.querySelectorAll('path[d^="m14 17 4 4"], path[d^="M18 13v8"]');
  animate(arrow, { translateY: [-3, 1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
