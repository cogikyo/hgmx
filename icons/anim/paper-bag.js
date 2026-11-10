import { animate, τmax } from "../motion.js";

export default function (el) {
  const bag = el.querySelectorAll("path");
  animate(bag, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
