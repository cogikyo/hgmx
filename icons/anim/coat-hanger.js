import { animate, τmax } from "../motion.js";

export default function (el) {
  const hanger = el.querySelector("g");
  animate(hanger, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
