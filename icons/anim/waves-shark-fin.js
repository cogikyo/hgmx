import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const waves = ["M2 15", "M2 21"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: 0.3, delay: stagger(0.06) });
  const fin = el.querySelector('path[d^="M17.3"]');
  animate(fin, { translateX: [-3, 1, 0] }, { duration: τmax });
}
