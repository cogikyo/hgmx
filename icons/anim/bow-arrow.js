import { animate, τmax } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [0, -1.5, 3, 0], translateY: [0, 1.5, -3, 0] }, { duration: τmax });
}
