import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelectorAll('path[d^="M13 2"], path[d^="M20 15.2"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const key = el.querySelectorAll('path[d="M17 2v6"], path[d="M17 4h2"], circle');
  animate(key, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
