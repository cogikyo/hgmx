import { animate, τmax } from "../motion.js";

export default function (el) {
  const croissant = el.querySelector("g");
  animate(croissant, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
