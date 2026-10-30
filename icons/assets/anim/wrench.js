import { animate, τmax } from "../motion.js";

export default function (el) {
  const wrench = el.querySelector("g");
  animate(wrench, { rotate: [0, -20, 8, 0] }, { duration: τmax });
}
