import { animate, stagger, τmax } from "../motion.js";

export default function (el) {
  const left = ["M9 5V3", "M9 10V9", "M9 15v-1", "M9 21v-2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(left, { translateX: [-6, 1, 0] }, { duration: τmax, delay: stagger(0.03) });

  const right = ["M15 5V3", "M15 10V9", "M15 15v-1", "M15 21v-2"].map((d) => el.querySelector(`path[d="${d}"]`));
  animate(right, { translateX: [6, -1, 0] }, { duration: τmax, delay: stagger(0.03) });
}
