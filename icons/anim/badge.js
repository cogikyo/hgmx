import { animate, τmax } from "../motion.js";

export default function (el) {
  const badge = el.querySelector("path");
  animate(badge, { scale: [1, 0.85, 1.1, 1] }, { duration: τmax });
}
