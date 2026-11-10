import { animate, τmax } from "../motion.js";

export default function (el) {
  const stem = el.querySelector("path");
  animate(stem, { scaleY: [0, 1.25, 1], opacity: [0, 1] }, { duration: τmax });
}
