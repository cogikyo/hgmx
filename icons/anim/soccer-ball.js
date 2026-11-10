import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector("g");
  animate(ball, { translateX: [0, 3, -1, 0], rotate: [0, 360] }, { duration: τmax });
}
