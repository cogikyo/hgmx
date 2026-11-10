import { animate, τmax } from "../motion.js";

export default function (el) {
  const pointer = el.querySelector("path");
  animate(pointer, { scale: [1, 0.85, 1.05, 1] }, { duration: τmax });
}
