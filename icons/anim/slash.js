import { animate, τmax } from "../motion.js";

export default function (el) {
  const slash = el.querySelector("path");
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
