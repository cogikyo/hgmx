import { animate, τmax } from "../motion.js";

export default function (el) {
  const screw = el.querySelectorAll("path");
  animate(screw, { translateX: [0, -2.5, 1, 0], translateY: [0, 2.5, -1, 0] }, { duration: τmax });
}
