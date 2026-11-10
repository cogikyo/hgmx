import { animate, τmax } from "../motion.js";

export default function (el) {
  const lotus = el.querySelector("g");
  animate(lotus, { scaleX: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
