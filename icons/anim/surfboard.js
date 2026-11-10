import { animate, τmax } from "../motion.js";

export default function (el) {
  const board = el.querySelector("g");
  animate(board, { translateY: [0, -2, 1, 0], rotate: [0, -6, 4, 0] }, { duration: τmax });
}
