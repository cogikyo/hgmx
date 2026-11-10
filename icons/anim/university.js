import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const walls = el.querySelector('path[d^="M22 7"]');
  animate(walls, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const windows = ["M6 12", "M18 12", "M6 16", "M18 16"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(windows, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
