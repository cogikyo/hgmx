import { animate, τmax } from "../motion.js";

export default function (el) {
  const ball = el.querySelector("g");
  animate(ball, { rotate: [0, -20, 10, 0] }, { duration: τmax });
}
