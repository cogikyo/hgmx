import { animate, τmax } from "../motion.js";

export default function (el) {
  const shell = el.querySelector("path");
  animate(shell, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
