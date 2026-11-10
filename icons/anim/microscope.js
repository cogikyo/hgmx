import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const arm = el.querySelector('path[d^="M14 22"]');
  animate(arm, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tube = el.querySelectorAll('path[d^="M9 12"], path[d^="M12 6V3"], path[d^="M9 14"]');
  animate(tube, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
