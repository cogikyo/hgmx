import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const flame = el.querySelector("path");
  animate(flame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(flame, { scaleY: [1, 0.85, 1.15, 1] }, { duration: τmax });
}
