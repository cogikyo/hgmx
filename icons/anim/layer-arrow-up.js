import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const layer = el.querySelector('path[d^="M7.674"]');
  animate(layer, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d^="M12 14"], path[d^="m9 7"]');
  animate(arrow, { translateY: [0, -3, 1, 0] }, { duration: τmid, delay: τmin });
}
