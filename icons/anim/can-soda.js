import { animate, τmax } from "../motion.js";

export default function (el) {
  const can = el.querySelector("g");
  animate(can, { rotate: [0, -8, 8, -4, 0] }, { duration: τmax });
}
