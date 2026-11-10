import { animate, stagger } from "../motion.js";

export default function (el) {
  const bands = ["M10 17", "M6 17", "M22 17"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(bands, { scale: [0, 1.25, 1], opacity: [0, 1] }, { duration: 0.3, delay: stagger(0.09) });
}
