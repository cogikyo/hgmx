import { animate, τmax } from "../motion.js";

export default function (el) {
  const curtains = el.querySelectorAll("g");
  animate(curtains, { scaleX: [1, 0.6, 1.05, 1] }, { duration: τmax });
}
