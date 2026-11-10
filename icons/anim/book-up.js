import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d="M12 13V7"], path[d^="m9 10"]');
  animate(arrow, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
