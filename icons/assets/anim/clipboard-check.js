import { animate, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector('path[d^="M16 4"]');
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const tick = el.querySelector('path[d^="m9 14"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
