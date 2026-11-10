import { animate, stagger, τmid, τmin } from "../motion.js";

export default function (el) {
  const rays = ["M9 7V5", "m13 8", "M13.95", "M9 19v-2", "M4.05", "M4 12H2", "m5 8"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.03) });

  const beam = el.querySelector('path[d^="M22 12"]');
  animate(beam, { scaleX: [0, 1] }, { duration: τmid, delay: τmin });
}
