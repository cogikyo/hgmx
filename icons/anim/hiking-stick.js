import { animate, τmax } from "../motion.js";

export default function (el) {
  const stick = el.querySelector("g");
  animate(stick, { rotate: [0, 10, -4, 0] }, { duration: τmax });
}
