import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = el.querySelector('path[d^="M20 19.645"]');
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M14.364"]');
  animate(pen, { rotate: [0, -12, 8, 0], translateX: [0, 3, -1, 0] }, { duration: τmax });

  const line = el.querySelector('path[d="M8 18h1"]');
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
