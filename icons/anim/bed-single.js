import { animate, τmax } from "../motion.js";

export default function (el) {
  const pillow = el.querySelector("path[transform-origin]");
  animate(pillow, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
