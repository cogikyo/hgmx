import { animate, τmax } from "../motion.js";

export default function (el) {
  const cherries = el.querySelector("g");
  animate(cherries, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
