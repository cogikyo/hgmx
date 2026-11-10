import { animate, τmax } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("path");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const pen = el.querySelector("g");
  animate(pen, { rotate: [0, -12, 8, 0], translateX: [0, 3, -1, 0] }, { duration: τmax });
}
