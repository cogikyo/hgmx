import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stem = el.querySelector('path[d^="M9 17c"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d^="M9 11.2"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
