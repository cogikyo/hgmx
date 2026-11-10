import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelectorAll("path");
  animate(ball, { translateX: [0, 3, -1, 0], translateY: [0, -3, 1, 0] }, { duration: τmax });
}
