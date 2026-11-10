import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector('path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const minus = el.querySelector('path[d="M9 14h6"]');
  animate(minus, { scale: [0, 1.25, 1], rotate: [0, 360] }, { duration: τmid, delay: τmin });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
