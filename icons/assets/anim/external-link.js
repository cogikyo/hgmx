import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [-4, 2, 0], translateY: [4, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const frame = el.querySelector('path[d^="M18 13"]');
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
}
