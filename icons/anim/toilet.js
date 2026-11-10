import { animate, τmax } from "../motion.js";

export default function (el) {
  const bowl = el.querySelector('path[d^="M7 12h13"]');
  animate(bowl, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
