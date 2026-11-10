import { animate, τmax } from "../motion.js";

export default function (el) {
  const paddle = el.querySelector("g");
  animate(paddle, { rotate: [0, -18, 10, -4, 0] }, { duration: τmax });
}
