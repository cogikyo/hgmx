import { animate, τmax } from "../motion.js";

export default function (el) {
  const tube = el.querySelector("g");
  animate(tube, { scaleX: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
