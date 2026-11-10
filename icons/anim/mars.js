import { animate, τmax } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { scale: [1, 1.25, 0.95, 1] }, { duration: τmax });
}
