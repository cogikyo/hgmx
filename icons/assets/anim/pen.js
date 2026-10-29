import { animate, τmax } from "../motion.js";

export default function (el) {
  const pen = el.querySelector("g");
  animate(pen, { rotate: [0, -12, 8, 0], translateX: [0, 3, -1, 0] }, { duration: τmax });
}
