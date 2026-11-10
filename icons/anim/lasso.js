import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const loop = el.querySelector('path[d^="M3.704"]');
  animate(loop, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });

  const knot = el.querySelector("circle");
  animate(knot, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
