import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const top = el.querySelector('path[d^="M3 5a"]');
  animate(top, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const frame = ["M21 9v", "M21 14v", "M21 19a", "M14 21h", "M9 21h", "M5 21a", "M3 14v", "M3 9v"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03, { startDelay: τmin }) });
}
