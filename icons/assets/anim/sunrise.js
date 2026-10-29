import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const arrow = el.querySelector("g");
  animate(arrow, { translateY: [4, -2, 0], opacity: [0, 1] }, { duration: τmax });

  const sun = el.querySelector('path[d^="M16 18"]');
  animate(sun, { translateY: [4, -1, 0] }, { duration: τmid });

  const rays = ["M2 18", "m4.93", "m19.07", "M20 18"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
