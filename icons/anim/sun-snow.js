import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const core = el.querySelector('path[d^="M10 9"]');
  animate(core, { scale: [0.5, 1.25, 1] }, { duration: τmax });

  const rays = ["M10 4", "m4.34", "M2 12", "m3.64", "M10 21"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });

  const snow = el.querySelector("g");
  animate(snow, { rotate: [0, 15, -5, 0] }, { duration: τmid, delay: τmin });
}
