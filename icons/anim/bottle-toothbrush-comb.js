import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M2 3a1"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const teeth = ["M18 6", "M18 10", "M18 14", "M18 18"].map((d) => el.querySelector(`path[d="${d}h4"]`));
  animate(teeth, { scaleX: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
