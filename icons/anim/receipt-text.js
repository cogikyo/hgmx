import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const paper = el.querySelector('path[d^="M4 3"]');
  animate(paper, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const lines = ["M14 8H8", "M16 12H8", "M13 16H8"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { opacity: [0, 1], x: [-3, 0] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
