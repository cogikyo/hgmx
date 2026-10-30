import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const ridge = el.querySelector("path");
  animate(ridge, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmax });

  const mountain = el.querySelector("g");
  animate(mountain, { scaleY: [0.5, 1.1, 1] }, { duration: τmid });
}
