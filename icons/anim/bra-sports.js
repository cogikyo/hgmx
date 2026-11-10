import { animate, τmax } from "../motion.js";

export default function (el) {
  const bra = el.querySelector("g");
  animate(bra, { rotate: [0, -8, 5, -2, 0] }, { duration: τmax });
}
