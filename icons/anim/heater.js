import { animate, stagger } from "../motion.js";

export default function (el) {
  const heat = ["M11 8", "M15.5 8"].map((d) => el.querySelector(`path[d^="${d}c"]`));
  animate(heat, { pathLength: [0, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.12) });
}
