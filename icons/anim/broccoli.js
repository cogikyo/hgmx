import { animate, τmax } from "../motion.js";

export default function (el) {
  const broccoli = el.querySelector("g");
  animate(broccoli, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
