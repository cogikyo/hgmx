import { animate, τmax } from "../motion.js";

export default function (el) {
  const body = el.querySelector("path");
  animate(body, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });
}
