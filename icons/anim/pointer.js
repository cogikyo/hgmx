import { animate, τmax } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { scale: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
