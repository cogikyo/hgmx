import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector("g");
  animate(ball, { translateY: [0, -2, 1, 0], rotate: [0, -15, 6, 0] }, { duration: τmax });
}
