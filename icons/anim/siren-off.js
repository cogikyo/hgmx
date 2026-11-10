import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const rays = ["M2 12h", "M12 2v", "M18.5 ", "M21 12h"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const slash = el.querySelector('path[d^="m2 2"]');
  animate(slash, { pathLength: [0, 1], opacity: [0, 1] }, { duration: τmid, delay: τmin });
}
