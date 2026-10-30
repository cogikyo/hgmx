import { animate, stagger, τmid } from "../motion.js";

export default function (el) {
  const ruler = el.querySelector('path[d^="M21.3"]');
  animate(ruler, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ticks = ["m8.5", "m11.5", "m14.5", "m17.5"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(ticks, { scale: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
