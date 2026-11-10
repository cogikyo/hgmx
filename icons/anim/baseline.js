import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const rule = el.querySelector('path[d="M4 20h16"]');
  animate(rule, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const letter = el.querySelectorAll('path[d^="m6 16"], path[d="M8 12h8"]');
  animate(letter, { translateY: [-3, 1, 0] }, { duration: τmid, delay: τmin });
}
