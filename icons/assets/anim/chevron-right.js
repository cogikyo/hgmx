import { animate, τmax } from "../motion.js";

export default function (el) {
  const chevron = el.querySelector("path");
  animate(chevron, { translateX: [0, 3, -1, 0] }, { duration: τmax });
}
