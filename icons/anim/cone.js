import { animate, τmax } from "../motion.js";

export default function (el) {
  const sides = el.querySelector("path");
  animate(sides, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
