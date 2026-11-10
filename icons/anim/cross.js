import { animate, τmax } from "../motion.js";

export default function (el) {
  const cross = el.querySelector("path");
  animate(cross, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
