import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const capsule = el.querySelector("g");
  animate(capsule, { scale: [1, 0.85, 1] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
