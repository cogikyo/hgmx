import { animate, τmax } from "../motion.js";

export default function (el) {
  const beam = el.querySelector('path[d^="M3 7"]');
  animate(beam, { rotate: [0, -8, 5, 0] }, { duration: τmax });

  const left = el.querySelector('path[d^="m5 8"]');
  animate(left, { translateY: [0, 1, -0.6, 0] }, { duration: τmax });

  const right = el.querySelector('path[d^="m19 8"]');
  animate(right, { translateY: [0, -1, 0.6, 0] }, { duration: τmax });
}
