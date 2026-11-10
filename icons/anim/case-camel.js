import { animate, τmax } from "../motion.js";

export default function (el) {
  const type = el.querySelector("g");
  animate(type, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
