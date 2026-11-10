import { animate, τmax } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("path");
  animate(lens, { scaleX: [1, 0.8, 1.1, 1] }, { duration: τmax });
}
