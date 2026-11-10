import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const hull = ["M22 7.7", "M10 21.9", "m10 14"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(hull, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const ribs = [el.querySelector('path[d^="M14 19.8"]'), el.querySelector('path[d^="M18 17.5"]')];
  animate(ribs, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmin, delay: stagger(0.12, { startDelay: τmin }) });
}
