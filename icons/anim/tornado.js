import { animate, stagger } from "../motion.js";

export default function (el) {
  const bands = ["M21 4", "M18 8", "M19 12", "M16 16", "M11 20"].map((d) => el.querySelector(`path[d^="${d}"]`));
  animate(bands, { translateX: [0, -3, 2, 0] }, { duration: 0.3, delay: stagger(0.06) });
}
