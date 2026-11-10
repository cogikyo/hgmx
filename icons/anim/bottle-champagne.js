import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const bottle = el.querySelector('path[d^="M9.7 21.3"]');
  animate(bottle, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });

  const fizz = ["M12 3", "M8 2", "M21 12", "M22 16"].map((d) => el.querySelector(`path[d="${d}h.01"]`));
  animate(fizz, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
