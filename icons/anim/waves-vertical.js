import { animate, stagger } from "../motion.js";

export default function (el) {
  const waves = ["M5 2", "M12 2", "M19 2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(waves, { pathLength: [0, 1], opacity: [0, 1], translateY: [-3, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
