import { animate, τmax } from "../motion.js";

export default function (el) {
  const petals = el.querySelector("path[transform-origin]");
  animate(petals, { scale: [0.6, 1.15, 1] }, { duration: τmax });
}
