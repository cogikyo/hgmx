import { animate, τmax } from "../motion.js";

export default function (el) {
  const tire = el.querySelector("g");
  animate(tire, { scaleY: [1, 0.88, 1.04, 1], scaleX: [1, 1.06, 0.98, 1] }, { duration: τmax });
}
