import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const pen = el.querySelector("g");
  animate(pen, { rotate: [0, -6, 3, 0], translateX: [0, 2.5, -1, 0] }, { duration: τmid });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
