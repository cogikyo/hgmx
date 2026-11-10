import { animate, τmax } from "../motion.js";

export default function (el) {
  const board = el.querySelector("g");
  animate(board, { translateX: [0, -3, 1, 0], translateY: [0, 3, -1, 0], rotate: [0, 6, -3, 0] }, { duration: τmax });
}
