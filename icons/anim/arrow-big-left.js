import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("path");
  animate(arrow, { translateX: [0, -4, 1, 0] }, { duration: τmax });
}
