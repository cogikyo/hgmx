import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="m7 4"]');
  animate(top, { translateY: [0, 1.5, -0.5, 0] }, { duration: τmax });

  const bottom = el.querySelector('path[d^="m7 20"]');
  animate(bottom, { translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });
}
