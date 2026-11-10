import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const letter = el.querySelectorAll('path[d^="m2 16"], path[d^="M3.304"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const arrow = el.querySelectorAll('path[d^="m14 12"], path[d="M18 16V7"]');
  animate(arrow, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
