import { animate, stagger } from "../motion.js";

export default function (el) {
  const blocks = ["M6.5 3.65", "M10.1 5.8", "M15.9 10.4", "M3.3 7.1", "M2 12c", "M17.7 15.1"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(blocks, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
