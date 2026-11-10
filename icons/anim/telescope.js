import { animate, τmax } from "../motion.js";

export default function (el) {
  const tube = el.querySelector("g");
  animate(tube, { rotate: [0, -10, 6, -2, 0] }, { duration: τmax });
}
