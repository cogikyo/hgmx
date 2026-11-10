import { animate, τmax } from "../motion.js";

export default function (el) {
  const burger = el.querySelector("g");
  animate(burger, { scaleY: [1, 0.8, 1.05, 1], scaleX: [1, 1.08, 0.98, 1] }, { duration: τmax });
}
