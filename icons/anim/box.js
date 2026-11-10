import { animate, τmax } from "../motion.js";

export default function (el) {
  const box = el.querySelector("g");
  animate(box, { translateY: [0, -4, 1, 0], scaleY: [1, 1.1, 0.9, 1] }, { duration: τmax });
}
