import { animate, τmax } from "../motion.js";

export default function (el) {
  const plant = el.querySelector("path[transform-origin]");
  animate(plant, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
