import { animate, τmax } from "../motion.js";

export default function (el) {
  const jacket = el.querySelector("path");
  animate(jacket, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
