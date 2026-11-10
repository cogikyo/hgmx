import { animate, τmax } from "../motion.js";

export default function (el) {
  const bag = el.querySelector("g");
  animate(bag, { translateY: [0, -2, 0], rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
}
