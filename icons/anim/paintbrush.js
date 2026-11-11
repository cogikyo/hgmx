import { animate, τmax } from "../motion.js";

export default function (el) {
  const brush = el.querySelector("g");
  animate(brush, { rotate: [0, -8, 8, -4, 0], translateX: [0, -2, 2, -1, 0] }, { duration: τmax });
}
