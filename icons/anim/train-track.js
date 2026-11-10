import { animate, stagger } from "../motion.js";

export default function (el) {
  const ties = ["m2 14", "m5 11", "m8 8", "m11 5", "m14 2"].map((d) => el.querySelector(`path[d^="${d} "]`));
  animate(ties, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
