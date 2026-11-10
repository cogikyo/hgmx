import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const fluid = el.querySelector('path[d^="M6 11"]');
  animate(fluid, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });

  const line = el.querySelector('path[d^="M12 18"]');
  animate(line, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
