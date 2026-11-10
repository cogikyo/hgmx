import { animate, τmax } from "../motion.js";

export default function (el) {
  const up = el.querySelectorAll('path[d^="M12 8"], path[d^="m15 5"]');
  animate(up, { translateY: [4, -1, 0], opacity: [0, 1] }, { duration: τmax });

  const down = el.querySelectorAll('path[d^="M12 22"], path[d^="m15 19"]');
  animate(down, { translateY: [-4, 1, 0], opacity: [0, 1] }, { duration: τmax });
}
