import { animate, τmax } from "../motion.js";

export default function (el) {
  const germ = el.querySelector("g");
  animate(germ, { scaleX: [1, 1.1, 0.92, 1], scaleY: [1, 0.9, 1.08, 1] }, { duration: τmax });
}
