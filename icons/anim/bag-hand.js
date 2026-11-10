import { animate, τmax } from "../motion.js";

export default function (el) {
  const bag = el.querySelector("g");
  animate(bag, { rotate: [0, -10, 8, -4, 0] }, { duration: τmax });
}
