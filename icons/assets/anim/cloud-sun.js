import { animate, stagger, τmax, τmid, τmin } from "../motion.js";

export default function (el) {
  const cloud = el.querySelector('path[d^="M13 22"]');
  animate(cloud, { translateX: [-3, 1, 0], opacity: [0.5, 1] }, { duration: τmid });

  const sun = el.querySelector('path[d^="M15.947"]');
  animate(sun, { scale: [0.5, 1.25, 1] }, { duration: τmax });

  const rays = ["m4.93", "M12 2", "m19.07", "M20 12"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03, { startDelay: τmin }) });
}
