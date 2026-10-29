import { animate, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { rotate: [0, -20, 2, 0] }, { duration: τmax });

  const board = el.querySelector(":scope > path");
  animate(board, { translateY: [0, 1, 0] }, { duration: τmin, delay: τmid });
}
