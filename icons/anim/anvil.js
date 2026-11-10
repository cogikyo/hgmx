import { animate, τmax } from "../motion.js";

export default function (el) {
  const anvil = el.querySelector("g");
  animate(anvil, { scaleY: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
