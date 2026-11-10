import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const book = el.querySelector('path[d^="M4 19.5"]');
  animate(book, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const phones = el.querySelector("g");
  animate(phones, { scaleX: [1, 0.75, 1.15, 1] }, { duration: τmid, delay: τmin });
}
