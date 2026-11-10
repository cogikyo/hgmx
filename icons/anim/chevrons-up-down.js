import { animate, τmax } from "../motion.js";

export default function (el) {
  const up = el.querySelector('path[d^="m7 9"]');
  animate(up, { translateY: [0, -3, 1, 0] }, { duration: τmax });

  const down = el.querySelector('path[d^="m7 15"]');
  animate(down, { translateY: [0, 3, -1, 0] }, { duration: τmax });
}
