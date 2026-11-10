import { animate, stagger } from "../motion.js";

export default function (el) {
  const rays = ["M4.208 4.5", "M3 9H2", "m4.203 13.5", "m19.797 13.5", "M21 9h1", "m19.792 4.5"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(rays, { scale: [0, 1.75, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
