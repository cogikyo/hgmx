import { animate, τmax } from "../motion.js";

export default function (el) {
  const top = el.querySelector("path");
  animate(top, { scaleY: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
