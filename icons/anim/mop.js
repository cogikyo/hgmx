import { animate, τmax } from "../motion.js";

export default function (el) {
  const mop = el.querySelector("g");
  animate(mop, { translateX: [0, -2.5, 1.5, 0] }, { duration: τmax });
}
