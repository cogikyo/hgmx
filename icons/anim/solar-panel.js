import { animate, stagger } from "../motion.js";

export default function (el) {
  const rays = ["M11 2h2", "m8.66 7.66", "M3 10v2"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
