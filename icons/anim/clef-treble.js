import { animate, τmax } from "../motion.js";

export default function (el) {
  const clef = el.querySelector("path");
  animate(clef, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
