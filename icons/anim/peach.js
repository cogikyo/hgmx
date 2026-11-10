import { animate, τmax } from "../motion.js";

export default function (el) {
  const peach = el.querySelector("g");
  animate(peach, { scaleY: [1, 0.85, 1.1, 1], scaleX: [1, 1.1, 0.95, 1] }, { duration: τmax });
}
