import { animate, τmax } from "../motion.js";

export default function (el) {
  const board = el.querySelector("g");
  animate(board, { scaleY: [1, 0.9, 1.05, 1] }, { duration: τmax });
}
