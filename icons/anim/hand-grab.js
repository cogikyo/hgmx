import { animate, τmax } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { scaleY: [1, 0.75, 1.05, 1] }, { duration: τmax });
}
