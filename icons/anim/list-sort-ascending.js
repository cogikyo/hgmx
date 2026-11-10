import { animate, stagger } from "../motion.js";

export default function (el) {
  const lines = ["M9 5H3", "M15 12H3", "M3 19h18"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(lines, { scaleX: [0, 1.15, 1] }, { duration: 0.3, delay: stagger(0.06) });
}
