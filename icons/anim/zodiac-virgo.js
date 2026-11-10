import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const strokes = ["M6 19", "M6 5.5", "M11 5.5"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(strokes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06) });

  const loop = el.querySelector('path[d^="M16 11.5"]');
  animate(loop, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
