import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const board = el.querySelector("path");
  animate(board, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const clip = el.querySelector("rect");
  animate(clip, { translateY: [-3, 1, 0] }, { duration: τmin, delay: τmid });
}
