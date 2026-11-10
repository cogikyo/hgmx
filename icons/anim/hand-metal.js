import { animate, τmax } from "../motion.js";

export default function (el) {
  const hand = el.querySelector("g");
  animate(hand, { translateY: [0, -3, 1, 0], rotate: [0, 10, -4, 0] }, { duration: τmax });
}
