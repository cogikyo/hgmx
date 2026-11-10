import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("path");
  animate(arrow, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
