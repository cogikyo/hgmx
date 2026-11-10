import { animate, τmax } from "../motion.js";

export default function (el) {
  const mask = el.querySelector("g");
  animate(mask, { translateY: [0, -3, 1, 0], rotate: [0, 8, -4, 0] }, { duration: τmax });
}
