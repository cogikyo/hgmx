import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const rule = el.querySelector('path[d="M7 17h10"]');
  animate(rule, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });
  const letter = el.querySelectorAll('path[d^="m9 13"], path[d="M9.5 12h5"]');
  animate(letter, { translateY: [-2, 0.5, 0] }, { duration: τmid, delay: τmin });
}
