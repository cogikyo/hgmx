import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const cup = el.querySelector('path[d^="M16 8"]');
  animate(cup, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const steam = ["M6 2", "M10 2", "M14 2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(steam, { translateY: [3, -1, 0], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
