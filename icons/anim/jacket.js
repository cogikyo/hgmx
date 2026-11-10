import { animate, τmax } from "../motion.js";

export default function (el) {
  const jacket = el.querySelector("g");
  animate(jacket, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
