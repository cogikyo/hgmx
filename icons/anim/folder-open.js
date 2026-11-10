import { animate, τmax } from "../motion.js";

export default function (el) {
  const folder = el.querySelector("path");
  animate(folder, { scaleY: [1, 1.15, 0.95, 1], translateY: [0, -1, 0] }, { duration: τmax });
}
