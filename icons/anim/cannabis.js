import { animate, τmax } from "../motion.js";

export default function (el) {
  const leaf = el.querySelector("g");
  animate(leaf, { rotate: [0, -10, 6, -3, 0] }, { duration: τmax });
}
