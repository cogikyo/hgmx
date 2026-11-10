import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const compass = el.querySelector("g");
  animate(compass, { rotate: [0, -10, 6, 0] }, { duration: τmax });

  const arc = el.querySelector('path[d^="M19.136"]');
  animate(arc, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
