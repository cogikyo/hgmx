import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector("g");
  animate(ball, { translateY: [0, -4, 1, 0], scaleX: [1, 0.9, 1.1, 1], scaleY: [1, 1.1, 0.9, 1] }, { duration: τmax });
}
