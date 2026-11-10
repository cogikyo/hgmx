import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const arc = el.querySelector('path[d^="M15.6"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  const needle = el.querySelector('path[d^="M13.4"]');
  animate(needle, { rotate: [-90, 20, -10, 0] }, { duration: τmax });
}
