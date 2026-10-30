import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ribbon = el.querySelector('path[d^="M10 2"]');
  animate(ribbon, { scaleY: [0, 1.25, 1] }, { duration: τmid, delay: τmin });
}
