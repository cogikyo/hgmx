import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const frame = el.querySelectorAll('path[d^="M10 3"], path[d^="M10 21"], path[d^="M22 21"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const door = el.querySelector("g");
  animate(door, { scaleX: [0.3, 1.15, 1] }, { duration: τmax });
}
