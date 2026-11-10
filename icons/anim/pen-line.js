import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pen = el.querySelector("g");
  animate(pen, { rotate: [0, -12, 8, 0], translateX: [0, 3, -1, 0] }, { duration: τmax });

  const line = el.querySelector('path[d^="M13 21"]');
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
