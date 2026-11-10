import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const glass = el.querySelector('path[d^="M20 14.54"]');
  animate(glass, { pathLength: [0, 1], opacity: [0.5, 1] }, { duration: τmid });
  animate(glass, { scaleY: [1, 1.15, 0.95, 1] }, { duration: τmid, delay: τmin });
  const rays = ["M12 2", "m4.93", "M2 12", "m6.34"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
