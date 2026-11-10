import { animate, τmax } from "../motion.js";

export default function (el) {
  const blades = el.querySelector("g");
  animate(blades, { rotate: [0, 360], scale: [1, 0.85, 1] }, { duration: τmax });
  animate(blades.querySelectorAll("path"), { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmax });
}
