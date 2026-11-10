import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const sides = [el.querySelector('path[d^="M13 6"]'), el.querySelector('path[d^="M11 18"]')];
  animate(sides, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(τmin) });
}
