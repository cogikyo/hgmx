import { animate, τmax } from "../motion.js";

export default function (el) {
  const arc = el.querySelector("path");
  animate(arc, { rotate: [0, 720] }, { duration: τmax });
  animate(arc, { pathLength: [0.25, 1] }, { duration: τmax });
}
