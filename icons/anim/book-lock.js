import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelectorAll('path[d^="M20 15"], path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const shackle = el.querySelector('path[d^="M18 6"]');
  animate(shackle, { translateY: [0, 1.5, 0] }, { duration: τmid, delay: τmin });

  const body = el.querySelector("rect");
  animate(body, { scale: [1, 1.1, 1] }, { duration: τmin, delay: τmid });
}
