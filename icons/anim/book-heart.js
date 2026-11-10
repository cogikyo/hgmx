import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const heart = el.querySelector('path[d^="M8.62"]');
  animate(heart, { scale: [1, 1.25, 0.9, 1.15, 1] }, { duration: τmid, delay: τmin });
}
