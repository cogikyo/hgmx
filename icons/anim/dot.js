import { animate, τmax } from "../motion.js";

export default function (el) {
  const dot = el.querySelector("circle");
  animate(dot, { scale: [1, 1.75, 1] }, { duration: τmax });
}
