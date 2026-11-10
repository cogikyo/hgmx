import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const dashes = ["M13.77 3", "M20.18 6", "M20.957", "M20.18 17", "M13.771 20", "M6.26 20", "M3.043", "M6.26 3"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dashes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03) });
}
