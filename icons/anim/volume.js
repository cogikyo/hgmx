import { animate, τmax } from "../motion.js";

export default function (el) {
  const speaker = el.querySelector("path");
  animate(speaker, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
