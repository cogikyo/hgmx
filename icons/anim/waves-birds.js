import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const waves = ["M2 15", "M2 21"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: 0.3, delay: stagger(0.06) });
  const birds = [el.querySelector('path[d^="M2 3"]'), el.querySelector('path[d^="M10 8"]')];
  animate(birds, { translateY: [0, -2, 0.5, 0] }, { duration: 0.3, delay: stagger(0.12, { startDelay: τmin }) });
}
