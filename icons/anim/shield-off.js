import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const halves = el.querySelectorAll("g path");
  animate(halves, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  const shield = el.querySelector("g");
  animate(shield, { scale: [1, 0.85, 1] }, { duration: τmin, delay: τmid });
}
