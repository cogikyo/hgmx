import { animate, τmax } from "../motion.js";

export default function (el) {
  const bloom = el.querySelector("g");
  animate(bloom, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
