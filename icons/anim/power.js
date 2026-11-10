import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const stem = el.querySelector('path[d="M12 2v10"]');
  animate(stem, { translateY: [0, 3, -1, 0] }, { duration: τmax });

  const arc = el.querySelector('path[d^="M18.4"]');
  animate(arc, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
