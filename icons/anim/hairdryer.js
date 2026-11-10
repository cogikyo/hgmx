import { animate, τmax } from "../motion.js";

export default function (el) {
  const dryer = el.querySelector("g");
  animate(dryer, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
