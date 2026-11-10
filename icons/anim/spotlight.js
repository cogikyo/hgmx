import { animate, stagger } from "../motion.js";

export default function (el) {
  const rays = ["m19 12.5", "m17 16", "M15.295"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
