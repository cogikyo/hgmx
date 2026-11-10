import { animate, τmax } from "../motion.js";

export default function (el) {
  const magnet = el.querySelectorAll("path");
  animate(magnet, { translateX: [0, -2, 0.5, 0], translateY: [0, 2, -0.5, 0] }, { duration: τmax });
}
