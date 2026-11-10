import { animate, τmax } from "../motion.js";

export default function (el) {
  const arm = el.querySelector("g");
  animate(arm, { rotate: [0, -15, 8, 0] }, { duration: τmax });
}
