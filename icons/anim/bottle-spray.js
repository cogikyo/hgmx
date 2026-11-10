import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M10 8"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const trigger = el.querySelector('path[d^="M14 6c"]');
  animate(trigger, { rotate: [0, 20, 0] }, { duration: τmin });

  const mist = ["M20 5.5", "M22 2", "M22 9"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(mist, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06, { startDelay: τmin }) });
}
