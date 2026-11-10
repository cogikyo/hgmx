import { animate, τmax } from "../motion.js";

export default function (el) {
  const pretzel = el.querySelector("path");
  animate(pretzel, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
