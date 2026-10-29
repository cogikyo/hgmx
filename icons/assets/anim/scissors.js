import { animate, τmax } from "../motion.js";

export default function (el) {
  const [upper, lower] = el.querySelectorAll("g");
  animate(upper, { rotate: [0, 15, -5, 0] }, { duration: τmax });
  animate(lower, { rotate: [0, -15, 5, 0] }, { duration: τmax });
}
