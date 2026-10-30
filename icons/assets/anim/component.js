import { animate, stagger } from "../motion.js";

export default function (el) {
  const parts = ["M8.916 4.674", "M15.536", "M8.916 17.912", "M2.297"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(parts, { scale: [0, 1.5, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
