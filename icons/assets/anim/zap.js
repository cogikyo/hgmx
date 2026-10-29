import { animate, τmax } from "../motion.js";

export default function (el) {
  const bolt = el.querySelector("path");
  animate(bolt, { translateX: [3, -1, 0], translateY: [-4, 1, 0] }, { duration: τmax });
}
