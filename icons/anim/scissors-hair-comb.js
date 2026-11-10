import { animate, τmax } from "../motion.js";

export default function (el) {
  const [left, right] = el.querySelectorAll("g");
  animate(left, { rotate: [0, 15, -5, 0] }, { duration: τmax });
  animate(right, { rotate: [0, -15, 5, 0] }, { duration: τmax });
}
