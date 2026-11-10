import { animate, τmax } from "../motion.js";

export default function (el) {
  const pillow = el.querySelector("path");
  animate(pillow, { scaleY: [1, 0.8, 1.05, 1], scaleX: [1, 1.08, 1] }, { duration: τmax });
}
