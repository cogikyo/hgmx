import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const axes = el.querySelector('path[d^="M4 4v15"]');
  animate(axes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid });

  const depth = ["M4.293", "m9 15", "M13.5"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(depth, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
