import { animate, τmax } from "../motion.js";

export default function (el) {
  const loop = el.querySelector("path");
  animate(loop, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
