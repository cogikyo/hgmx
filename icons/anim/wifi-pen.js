import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const waves = ["M8.5", "M5", "M2 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });

  const pen = el.querySelector('path[d^="M21.378"]');
  animate(pen, { rotate: [-15, 15, -10, 0], translateX: [0, -3, 0] }, { duration: τmax });
}
