import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector("path");
  animate(bubble, { scaleY: [0, 1.1, 1] }, { duration: τmax });
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
