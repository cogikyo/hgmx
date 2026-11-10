import { animate, τmax } from "../motion.js";

export default function (el) {
  const cap = el.querySelector("g");
  animate(cap, { translateY: [0, -3, 0], rotate: [0, 12, -4, 0] }, { duration: τmax });
}
