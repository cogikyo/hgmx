import { animate, τmax } from "../motion.js";

export default function (el) {
  const gavel = el.querySelector("g");
  animate(gavel, { rotate: [0, -20, 2, -8, 0] }, { duration: τmax });
}
