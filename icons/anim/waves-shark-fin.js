import { animate, τmax } from "../motion.js";

export default function (el) {
  const swell = el.querySelector('path[d^="M2 21"]');
  animate(swell, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: 0.3 });
  const fin = el.querySelector('path[d^="M17.3"]');
  animate(fin, { translateX: [0, -3, 1, 0] }, { duration: τmax });
}
