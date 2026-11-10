import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M4 14.899"]');
  animate(cloud, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [6, -2, 0], opacity: [0, 1] }, { duration: τmax });
}
