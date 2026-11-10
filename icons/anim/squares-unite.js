import { animate, τmax } from "../motion.js";

export default function (el) {
  const outline = el.querySelector("path");
  animate(outline, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
