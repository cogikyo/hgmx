import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const folder = el.querySelector('path[d^="M20 20"]');
  animate(folder, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lens = el.querySelector("g");
  animate(lens, { scale: [1.35, 1] }, { duration: τmax });
  animate(lens, { translateX: [0, -1.5, 1, 0], translateY: [0, -1.5, 0.5, 0] }, { duration: τmax });
}
