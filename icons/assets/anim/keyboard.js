import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector("rect");
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const keys = el.querySelectorAll('path[d$="h.01"]');
  animate(keys, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const space = el.querySelector('path[d="M7 16h10"]');
  animate(space, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
