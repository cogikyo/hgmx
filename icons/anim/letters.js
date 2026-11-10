import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const glyphs = ["m8 10", "M3 18", "M21 15"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(glyphs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
  const bar = el.querySelector('path[d="M15 8H9"]');
  animate(bar, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: τmid });
}
