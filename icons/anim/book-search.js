import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const book = el.querySelectorAll('path[d^="M11 22"], path[d^="M3 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -4, 2, 0], translateY: [0, -3, -2, 0] }, { duration: τmax });
}
