import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const frame = el.querySelector("rect");
  animate(frame, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const letters = ["M10 15", "M10 13", "M14 14.5"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(letters, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.06, { startDelay: τmin }) });
}
