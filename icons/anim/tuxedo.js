import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const lapels = el.querySelector('path[d^="M18 3h1"]');
  animate(lapels, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
  const bow = el.querySelector('path[d^="M10 3v2"]');
  animate(bow, { rotate: [0, -20, 10, 0] }, { duration: τmid, delay: τmin });
}
