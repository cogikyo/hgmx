import { animate, τmax } from "../motion.js";

export default function (el) {
  const block = el.querySelector("rect");
  animate(block, { scaleX: [1, 0.7, 1.1, 1] }, { duration: τmax });
}
