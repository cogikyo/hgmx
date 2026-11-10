import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const ring = el.querySelector("circle");
  animate(ring, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 90] }, { duration: τmid });

  const letter = el.querySelector('path[d^="m8 16"]');
  animate(letter, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const bar = el.querySelector('path[d="M15 14H9"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
