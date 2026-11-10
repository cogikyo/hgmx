import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const dashes = ["M10.17", "m15.874", "m18.849", "M21.824", "M14 21", "M8 21", "M4.024", "m5.136", "m8.102"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dashes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03) });
}
