import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const grill = el.querySelector('path[d^="M3 8"]');
  animate(grill, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const smoke = ["M6 4", "M12 4", "M18 4"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(smoke, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
