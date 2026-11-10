import { animate, τmax } from "../motion.js";

export default function (el) {
  const lens = el.querySelector("path");
  animate(lens, { scaleX: [1, 1.25, 0.9, 1] }, { duration: τmax });
}
