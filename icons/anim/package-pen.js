import { animate, τmax, τmid } from "../motion.js";

export default function (el) {
  const box = ["M12 14.34", "M21 9.114", "M3.29 7", "m7.5 4.27"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(box, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pen = el.querySelector('path[d^="M16.511"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
