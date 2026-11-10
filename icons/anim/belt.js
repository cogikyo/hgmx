import { animate, τmax } from "../motion.js";

export default function (el) {
  const belt = el.querySelector("g");
  animate(belt, { scaleX: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
