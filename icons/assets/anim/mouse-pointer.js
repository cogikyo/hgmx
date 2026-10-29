import { animate, τmax } from "../motion.js";

export default function (el) {
  const pointer = el.querySelectorAll("path");
  animate(pointer, { translateX: [0, -3, 1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
}
