import { animate, τmax } from "../motion.js";

export default function (el) {
  const brush = el.querySelectorAll("path");
  animate(brush, { translateX: [0, 2.5, -2.5, 1.5, 0], translateY: [0, -2.25, 2.25, -1.35, 0] }, { duration: τmax });
}
