import { animate, τmax } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { rotate: [0, -15, 12, -6, 0] }, { duration: τmax });
}
