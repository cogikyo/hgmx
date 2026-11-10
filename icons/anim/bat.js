import { animate, τmax } from "../motion.js";

export default function (el) {
  const bat = el.querySelectorAll("path");
  animate(bat, { scaleX: [1, 0.6, 1.1, 1], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });
}
