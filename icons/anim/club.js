import { animate, τmax } from "../motion.js";

export default function (el) {
  const suit = el.querySelector("g");
  animate(suit, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
