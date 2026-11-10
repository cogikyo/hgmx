import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pencil = el.querySelector("g");
  animate(pencil, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });

  const line = el.querySelector('path[d^="M13 21"]');
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
