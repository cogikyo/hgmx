import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelectorAll('path[d^="M8 4"], path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelectorAll('path[d="M11 14h10"], path[d^="m17 18"]');
  animate(arrow, { translateX: [0, 3, -1, 0] }, { duration: τmid, delay: τmin });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
