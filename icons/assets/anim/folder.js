import { animate, τmax } from "../motion.js";

export default function (el) {
  const folder = el.querySelector("path");
  animate(folder, { scaleY: [1, 0.85, 1.05, 1], translateY: [0, 1, 0] }, { duration: τmax });
}
