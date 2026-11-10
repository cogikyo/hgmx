import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const letter = el.querySelectorAll('path[d^="m6 16"], path[d="M8 12h8"]');
  animate(letter, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const squiggle = el.querySelector('path[d^="M4 21"]');
  animate(squiggle, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
