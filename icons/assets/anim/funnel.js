import { animate, τmax } from "../motion.js";

export default function (el) {
  const funnel = el.querySelector("path");
  animate(funnel, { scaleX: [1, 0.8, 1.1, 1], scaleY: [1, 1.1, 0.95, 1] }, { duration: τmax });
}
