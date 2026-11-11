import { animate, τmax } from "../motion.js";

export default function (el) {
  const trousers = el.querySelector("g");
  animate(trousers, { scaleY: [1, 0.84, 1, 0.9, 1] }, { duration: τmax });
}
