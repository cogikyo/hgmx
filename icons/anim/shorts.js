import { animate, τmax } from "../motion.js";

export default function (el) {
  const shorts = el.querySelector("g");
  animate(shorts, { scaleY: [1, 0.8, 1.08, 1] }, { duration: τmax });
}
