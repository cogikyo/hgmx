import { animate, τmax } from "../motion.js";

export default function (el) {
  const fist = el.querySelector("g");
  animate(fist, { scale: [1, 1.15, 0.95, 1] }, { duration: τmax });
}
