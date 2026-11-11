import { animate, τmax } from "../motion.js";

export default function (el) {
  const bra = el.querySelector("g");
  animate(bra, { translateY: [0, -3, 1, 0] }, { duration: τmax });
}
