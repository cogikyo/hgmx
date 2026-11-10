import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const dashes = ["M10.1", "M17.6", "M21.8", "M20.2", "M13.9", "M6.", "M2.18", "M3.7"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dashes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03) });

  const tick = el.querySelector('path[d^="m16 9"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
