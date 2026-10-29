import { animate, τmax } from "../motion.js";

export default function (el) {
  const plane = el.querySelector("path");
  animate(plane, { translateX: [0, 4, -1, 0], translateY: [0, -4, 1, 0], rotate: [0, -10, 5, 0] }, { duration: τmax });
}
