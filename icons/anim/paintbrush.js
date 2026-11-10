import { animate, τmax } from "../motion.js";

export default function (el) {
  const brush = el.querySelector("g");
  animate(brush, { translateX: [0, -1.5, 0.5, 0], translateY: [0, 1.5, -0.5, 0], rotate: [0, -8, 3, 0] }, { duration: τmax });
}
