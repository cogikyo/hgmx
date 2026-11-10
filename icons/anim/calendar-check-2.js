import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const page = ["M 21 13", "M 21 5", "M 19 3", "M 5 3", "M 3 5", "M 3 19", "M5 21"].map((d) => el.querySelector(`path[d^="${d} "]`));
  animate(page, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmin, delay: stagger(0.03) });
  const tick = el.querySelector('path[d^="m16 19"]');
  animate(tick, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
  animate(tick, { scale: [1, 1.25, 1] }, { duration: τmin, delay: τmid });
}
