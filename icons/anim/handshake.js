import { animate, τmax } from "../motion.js";

export default function (el) {
  const hands = el.querySelectorAll("path");
  animate(hands, { translateY: [0, 2, -2, 1, 0] }, { duration: τmax });
}
