import { animate, τmax } from "../motion.js";

export default function (el) {
  const bag = el.querySelector("g");
  animate(bag, { translateX: [0, 2, -0.5, 0] }, { duration: τmax });
}
