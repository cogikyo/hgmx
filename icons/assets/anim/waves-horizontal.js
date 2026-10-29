import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = ["M2 5", "M2 12", "M2 19"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { pathLength: [0, 1], opacity: [0, 1], translateX: [-3, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
