import { animate, τmax } from "../motion.js";

export default function (el) {
  const swirl = el.querySelector("g");
  animate(swirl, { rotate: [0, 180] }, { duration: τmax });
}
