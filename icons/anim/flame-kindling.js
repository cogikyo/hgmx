import { animate, τmax } from "../motion.js";

export default function (el) {
  const flame = el.querySelector("path[transform-origin]");
  animate(flame, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
