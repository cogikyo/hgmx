import { animate, τmax } from "../motion.js";

export default function (el) {
  const sleeper = el.querySelector("g");
  animate(sleeper, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
