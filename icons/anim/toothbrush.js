import { animate, τmax } from "../motion.js";

export default function (el) {
  const brush = el.querySelectorAll("path");
  animate(brush, { translateX: [0, 1.5, -1.5, 0], translateY: [0, -1.35, 1.35, 0] }, { duration: τmax });
}
