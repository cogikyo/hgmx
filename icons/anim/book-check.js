import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m9 9.5"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
