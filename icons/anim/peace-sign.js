import { animate, τmax } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
