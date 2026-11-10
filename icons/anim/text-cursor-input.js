import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const field = el.querySelectorAll('path[d^="M13 8"], path[d^="M5 16"]');
  animate(field, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const caret = el.querySelector("g");
  animate(caret, { scaleY: [0, 1.15, 1] }, { duration: τmid });
  animate(caret, { opacity: [1, 0, 1] }, { duration: τmin, delay: τmid });
}
