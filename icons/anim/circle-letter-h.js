import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const stems = el.querySelectorAll('path[d$="V8"]');
  animate(stems, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const bar = el.querySelector('path[d="M9 12h6"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
