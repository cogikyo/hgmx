import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bubble = el.querySelector('path[d^="M3.086"]');
  animate(bubble, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = ["M8 8", "M8 12", "M8 16"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(lines, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
