import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const pencil = el.querySelector("g");
  animate(pencil, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });

  const ferrule = el.querySelector('path[d^="m15 5"]');
  animate(ferrule, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
