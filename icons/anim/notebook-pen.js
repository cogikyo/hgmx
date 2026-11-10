import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const cover = el.querySelector('path[d^="M13.4"]');
  animate(cover, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M21.378"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
