import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const dashes = ["M10.1", "M17.609", "M21.818", "M20.28", "M13.9", "m6.163", "M2.182", "M3.721"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dashes, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });
}
