import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const reticle = el.querySelector("g");
  animate(reticle, { scale: [1.2, 0.95, 1], opacity: [0.5, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2 "]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
