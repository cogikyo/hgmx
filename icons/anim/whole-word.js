import { animate, τmax } from "../motion.js";

export default function (el) {
  const bracket = el.querySelector('path[d^="M22 17"]');
  animate(bracket, { scaleX: [0, 1.15, 1] }, { duration: τmax });
}
