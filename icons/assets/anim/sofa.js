import { animate, τmax } from "../motion.js";

export default function (el) {
  const sofa = el.querySelector("g");
  animate(sofa, { scaleY: [1, 0.85, 1.05, 1], scaleX: [1, 1.05, 1] }, { duration: τmax });
}
