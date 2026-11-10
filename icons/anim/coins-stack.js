import { animate, τmax } from "../motion.js";

export default function (el) {
  const stack = el.querySelector("g");
  animate(stack, { translateY: [0, -3, 1, 0], scaleY: [1, 1.05, 0.9, 1] }, { duration: τmax });
}
