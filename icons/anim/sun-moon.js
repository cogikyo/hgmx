import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const moon = el.querySelector('path[d^="M14.837"]');
  animate(moon, { rotate: [0, -30, 10, 0], scale: [1, 0.85, 1] }, { duration: τmax });

  const rays = ["M12 2", "m19 5", "M20 12"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
