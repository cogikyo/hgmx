import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const screen = el.querySelector('path[d^="M13 3H4"]');
  animate(screen, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateX: [-3, 1, 0], translateY: [3, -1, 0], opacity: [0, 1] }, { duration: τmax });
}
