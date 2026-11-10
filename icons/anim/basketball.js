import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelectorAll("circle, path");
  animate(ball, { translateY: [0, 3, -1, 0] }, { duration: τmax });
}
