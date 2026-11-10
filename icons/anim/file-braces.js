import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M6 22"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const braces = el.querySelectorAll('path[d^="M10 12"], path[d^="M14 18"]');
  animate(braces, { scaleY: [0.5, 1.15, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
