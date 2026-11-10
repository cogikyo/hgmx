import { animate, τmax } from "../motion.js";

export default function (el) {
  const worm = el.querySelector("g");
  animate(worm, { scaleX: [1, 1.1, 0.95, 1], scaleY: [1, 0.9, 1.05, 1] }, { duration: τmax });
}
