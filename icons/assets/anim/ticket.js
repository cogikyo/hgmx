import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const ticket = el.querySelector('path[d^="M2 9"]');
  animate(ticket, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const holes = ["M13 5v2", "M13 11v2", "M13 17v2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(holes, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: τmid, delay: stagger(0.06, { startDelay: τmin }) });
}
