import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const strands = el.querySelectorAll('path[d^="M15 2"], path[d^="M2 15"], path[d^="M9 22"]');
  animate(strands, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const pairs = ["M3.109", "m6.5", "m7 18", "m10 16", "m14 8", "m16.5", "m17 6", "m20 9"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(pairs, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });
}
