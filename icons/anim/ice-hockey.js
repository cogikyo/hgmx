import { animate, τmax } from "../motion.js";

export default function (el) {
  const stick = el.querySelector("g");
  animate(stick, { rotate: [0, -15, 6, 0] }, { duration: τmax });
}
