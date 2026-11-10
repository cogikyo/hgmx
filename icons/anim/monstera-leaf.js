import { animate, τmax } from "../motion.js";

export default function (el) {
  const leaf = el.querySelector("g");
  animate(leaf, { rotate: [0, -12, 8, -4, 0] }, { duration: τmax });
}
