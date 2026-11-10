import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const stems = el.querySelectorAll('path[d$="V8"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });

  const bar = el.querySelector('path[d^="M9 12"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
