import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const ridge = el.querySelector('path[d^="m8 3"]');
  animate(ridge, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const mountain = el.querySelector("g");
  animate(mountain, { scaleY: [0.5, 1.1, 1] }, { duration: τmid });

  const snow = el.querySelector('path[d^="M4.14"]');
  animate(snow, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
