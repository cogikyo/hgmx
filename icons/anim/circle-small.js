import { animate, τmax } from "../motion.js";

export default function (el) {
  const circle = el.querySelector("circle");
  animate(circle, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
