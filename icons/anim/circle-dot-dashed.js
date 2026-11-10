import { animate, stagger, τmin } from "../motion.js";

export default function (el) {
  const dashes = ["M10.1", "M17.6", "M21.8", "M20.2", "M13.9", "M6.", "M2.18", "M3.7"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(dashes, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.03) });

  const dot = el.querySelector("circle");
  animate(dot, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: τmin });
}
