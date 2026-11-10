import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector("rect");
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const keys = ["M6", "M10", "M14", "M18"].map((x) => el.querySelector(`path[d="${x} 12v4"]`));
  animate(keys, { scaleY: [1, 1.5, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
