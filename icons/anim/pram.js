import { animate, τmax } from "../motion.js";

export default function (el) {
  const pram = el.querySelector("g");
  animate(pram, { rotate: [0, -5, 4, -2, 0] }, { duration: τmax });
}
