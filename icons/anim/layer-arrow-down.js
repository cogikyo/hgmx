import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const layer = el.querySelectorAll('path[d^="M22"]');
  animate(layer, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d^="M12 10"], path[d^="m9 17"]');
  animate(arrow, { translateY: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });
}
