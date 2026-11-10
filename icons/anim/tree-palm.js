import { animate, τmax } from "../motion.js";

export default function (el) {
  const fronds = el.querySelector("g");
  animate(fronds, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
