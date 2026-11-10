import { animate, τmax } from "../motion.js";

export default function (el) {
  const bloom = el.querySelector("g");
  animate(bloom, { scale: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
