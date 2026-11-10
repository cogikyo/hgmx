import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const platter = el.querySelector("circle");
  animate(platter, { pathLength: [0, 1], opacity: [0.5, 1], rotate: [270, 630] }, { duration: τmax });

  const arm = el.querySelector('path[d^="M18 8"]');
  animate(arm, { rotate: [0, 10, -3, 0] }, { duration: τmid, delay: τmin });
}
