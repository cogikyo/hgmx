import { animate, τmax } from "../motion.js";

export default function (el) {
  const flippers = el.querySelector("g");
  animate(flippers, { scaleY: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
