import { animate, τmax } from "../motion.js";

export default function (el) {
  const heart = el.querySelector("path");
  animate(heart, { scale: [1, 1.25, 0.9, 1.15, 1] }, { duration: τmax });
}
