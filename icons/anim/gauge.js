import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const dial = el.querySelector('path[d^="M3.34"]');
  animate(dial, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const needle = el.querySelector('path[d^="m12 14"]');
  animate(needle, { rotate: [-90, 20, -10, 0] }, { duration: τmax });
}
