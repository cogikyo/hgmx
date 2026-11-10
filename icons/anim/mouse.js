import { animate, τmax } from "../motion.js";

export default function (el) {
  const wheel = el.querySelector("path");
  animate(wheel, { translateY: [0, 3, -1, 0] }, { duration: τmax });
}
