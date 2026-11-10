import { animate, τmax } from "../motion.js";

export default function (el) {
  const shorts = el.querySelector("g");
  animate(shorts, { scaleX: [1, 1.08, 0.97, 1] }, { duration: τmax });
}
