import { animate, τmax } from "../motion.js";

export default function (el) {
  const blocks = el.querySelectorAll("rect");
  animate(blocks, { scaleX: [1, 0.7, 1] }, { duration: τmax });
}
