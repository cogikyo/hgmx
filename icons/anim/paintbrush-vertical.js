import { animate, τmax } from "../motion.js";

export default function (el) {
  const brush = el.querySelector("g");
  animate(brush, { rotate: [0, -10, 6, 0] }, { duration: τmax });
}
