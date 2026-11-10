import { animate, τmax, τmin } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { rotate: [0, 16, 0] }, { duration: τmax, times: [0, 0.35, 1] });

  const board = el.querySelector(":scope > path");
  animate(board, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmin });
}
