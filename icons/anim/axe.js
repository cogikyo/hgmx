import { animate, τmax } from "../motion.js";

export default function (el) {
  const axe = el.querySelectorAll("path");
  animate(axe, { translateY: [0, -3, 1.5, 0] }, { duration: τmax });
}
