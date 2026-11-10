import { animate, τmax } from "../motion.js";

export default function (el) {
  const squiggle = el.querySelector("path");
  animate(squiggle, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
