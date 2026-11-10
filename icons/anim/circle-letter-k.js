import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const stem = el.querySelector('path[d="M9 16V8"]');
  animate(stem, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const arms = el.querySelectorAll('path[d="m11 12 4 4"], path[d="m15 8-6 6"]');
  animate(arms, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
